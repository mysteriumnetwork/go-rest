package gin

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mysteriumnetwork/go-rest/paging"
)

const ContextKey = "page_request"

type Options struct {
	DefaultSize uint64
	MaxSize     uint64
}

func Middleware(c *gin.Context) {
	MiddlewareWith(c, Options{
		DefaultSize: paging.DefaultSize,
		MaxSize:     paging.MaxSize,
	})
}

func MiddlewareWith(c *gin.Context, opts Options) {
	pageRequest, err := parse(c, opts)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Set(ContextKey, pageRequest)
}

func Request(c *gin.Context) *paging.Request {
	val, exists := c.Get(ContextKey)
	if !exists {
		return &paging.Request{
			Page:     1,
			PageSize: paging.DefaultSize,
		}
	}
	pager, ok := val.(*paging.Request)
	if !ok {
		return &paging.Request{
			Page:     1,
			PageSize: paging.DefaultSize,
		}
	}
	return pager
}

func parse(c *gin.Context, opts Options) (*paging.Request, error) {
	size, err := pageSizeFromQuery(c, opts.DefaultSize, opts.MaxSize)
	if err != nil {
		return nil, fmt.Errorf("could not parse page size: %w", err)
	}
	page, err := pageFromQuery(c)
	if err != nil {
		return nil, fmt.Errorf("could not parse page number: %w", err)
	}
	return &paging.Request{
		Page:     page,
		PageSize: size,
	}, nil
}

func pageSizeFromQuery(c *gin.Context, defaultSize, maxSize uint64) (uint64, error) {
	pageSizeStr, ok := c.GetQuery("page_size")
	size := defaultSize
	if ok {
		var err error
		size, err = strconv.ParseUint(pageSizeStr, 10, 64)
		if err != nil {
			return 0, err
		}
		if size == 0 {
			return 0, errors.New("page size must be greater than 0")
		}
	}
	if size > maxSize {
		size = maxSize
	}
	return size, nil
}

func pageFromQuery(c *gin.Context) (uint64, error) {
	page := uint64(1)
	pageStr, ok := c.GetQuery("page")
	if ok {
		var err error
		page, err = strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			return page, err
		}
		if page == 0 {
			return 0, errors.New("page number must be greater than 0")
		}
	}
	return page, nil
}
