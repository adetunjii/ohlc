package api

import (
	"bufio"
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/adetunjii/ohlc/db/model"
	"github.com/gin-gonic/gin"
)

var (
	ErrReadingFormData = errors.New("error reading form data")
	ErrInvalidFileType = errors.New("invalid file type. Only csv files are supported")
	ErrFailedRetry     = errors.New("failed to insert data during retry")
)

// 5TB max file size
const MaxFileSize = 5000 << 40

func (s *Server) priceDataRouteGroup(r *gin.RouterGroup) {
	priceData := r.Group("/price-data")

	priceData.GET("/list", s.getAllPrices)
	priceData.POST("/create", s.createPrice)
	priceData.POST("/upload-price-list", s.uploadPriceList)

	// to be called by a CRON JOB
	priceData.GET("/retry-failed", s.retryFailedPriceListInsertion)
}

type getPricesQuery struct {
	Page   int32  `form:"page" binding:"required,min=1"`
	Size   int32  `form:"size" binding:"required,max=100"`
	Symbol string `form:"symbol"`
}

func (s *Server) getAllPrices(ctx *gin.Context) {
	var queryParams getPricesQuery
	err := ctx.ShouldBindQuery(&queryParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	limit := queryParams.Size
	offset := (queryParams.Page - 1) * queryParams.Size

	listPriceParams := model.ListPriceParams{
		Limit:  limit,
		Offset: offset,
		Symbol: strings.Trim(queryParams.Symbol, " "),
	}

	priceList, count, err := s.svc.ListPrices(ctx, listPriceParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Succesfully fetched price list",
		"data": gin.H{
			"page":       queryParams.Page,
			"size":       queryParams.Size,
			"totalItems": count,
			"content":    priceList,
		},
	})
}

type priceDataRequest struct {
	Symbol string `json:"symbol" binding:"required"`
	Open   string `json:"open" binding:"required"`
	High   string `json:"high" binding:"required"`
	Low    string `json:"low" binding:"required"`
	Close  string `json:"close" binding:"required"`
	Unix   int64  `json:"unix" binding:"required"`
}

func (s *Server) createPrice(ctx *gin.Context) {

	var payload priceDataRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	priceData := model.PriceData{
		Symbol: payload.Symbol,
		Open:   payload.Open,
		Close:  payload.Close,
		High:   payload.High,
		Low:    payload.Low,
		Unix:   payload.Unix,
	}

	if err := priceData.Validate(); err != nil {
		s.logger.Error("error parsing invalid data %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := s.svc.CreatePrice(ctx, priceData); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Completed successfully",
	})
}

func (s *Server) uploadPriceList(ctx *gin.Context) {

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, MaxFileSize)

	formData, err := ctx.MultipartForm()
	if err != nil {
		s.logger.Error("error reading multipart form data %v", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrReadingFormData))
		return
	}

	for _, fh := range formData.File {
		for _, header := range fh {
			fileExt := path.Ext(header.Filename)

			// only csv files are supported
			if fileExt != ".csv" {
				ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidFileType))
				return
			}

			file, err := header.Open()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(ErrReadingFormData))
				return
			}
			defer file.Close()

			// send the data chunks to the service layer for processing
			bufReader := bufio.NewReader(file)
			if err := s.svc.BulkInsertPrice(ctx, bufReader); err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "file uploaded successfully",
	})

}

func (s *Server) retryFailedPriceListInsertion(ctx *gin.Context) {

	if err := s.svc.RetryFromDeadQueue(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrFailedRetry))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully inserted all rows",
	})
}
