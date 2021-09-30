package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

type Product struct {
	Username    string    `form:"username" binding:"required"`
	Name        string    `form:"name" binding:"required"`
	Category    string    `form:"category" binding:"required"`
	Price       int       `form:"price" binding:"gte=0"`
	Description string    `form:"description"`
	CreatedAt   time.Time `form:"created_at"`
}

type productHandler struct {
	sync.RWMutex
	products map[string]Product
}

func newProductHandler() *productHandler {
	return &productHandler{
		products: make(map[string]Product),
	}
}

func (u *productHandler) Create(c *gin.Context) {
	u.Lock()
	defer u.Unlock()

	// 1. 参数解析
	var product Product
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// curl -XPOST -H"Content-Type: application/json" -d'{"username":"colin","name":"iphone12","category":"phone","price":8000,"description":"cannot afford"}' http://127.0.0.1:8080/v1/products
	// 2. 参数校验
	if _, ok := u.products[product.Name]; ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %s already exist", product.Name)})
		return
	}
	product.CreatedAt = time.Now()

	// 3. 逻辑处理
	u.products[product.Name] = product
	log.Printf("Register product %s success", product.Name)

	// 4. 返回结果
	c.JSON(http.StatusOK, product)
}

func (u *productHandler) Get(c *gin.Context) {
	u.Lock()
	defer u.Unlock()

	product, ok := u.products[c.Param("name")]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Errorf("can not found product %s", c.Param("name"))})
		return
	}

	c.JSON(http.StatusOK, product)
}

func router() http.Handler {
	router := gin.Default()
	productHandler := newProductHandler()
	// 路由分组、中间件、认证
	v1 := router.Group("/v1")
	{
		productv1 := v1.Group("/products")
		{
			// 路由匹配
			productv1.POST("", productHandler.Create)
			productv1.GET(":name", productHandler.Get)
		}
	}

	return router
}

func main() {
	var std = log.New(os.Stdout, "", log.LstdFlags)

	// 一进程多端口
	insecureServer := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	secureServer := &http.Server{
		Addr:         ":8443",
		Handler:      router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	eg, ctx := errgroup.WithContext(context.Background())

	go func() {
		eg.Go(func() error {
			err := insecureServer.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				std.Println("insecureServer:", err)
			}
			return err
		})

		eg.Go(func() error {
			time.Sleep(3 * time.Second)
			fmt.Println("3秒结束")
			return errors.New("定时抛出错误")
			err := secureServer.ListenAndServeTLS("../cert/apiserver.pem", "../cert/apiserver-key.pem")
			if err != nil && err != http.ErrServerClosed {
				std.Println("secureServer:", err)
			}
			return err
		})

		if err := eg.Wait(); err != nil {
			std.Println("asdfsf")

			std.Println(err)
		}
	}()
	quit := make(chan os.Signal)

	var cancel context.CancelFunc

	go func() {
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		std.Println("Shutting down server...")
	}()

	select {
	case <-ctx.Done():
	case <-quit:
		_, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	if err := insecureServer.Shutdown(ctx); err != nil {
		std.Fatal("insecureServer forced to shutdown:", err)
	}
	if err := secureServer.Shutdown(ctx); err != nil {
		std.Fatal("secureServer forced to shutdown:", err)
	}

	std.Println("Server exiting")
}
