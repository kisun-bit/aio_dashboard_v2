package core

import (
	"bytes"
	innerctx "context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kisun-bit/aio_dashboard/pkg/trace"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HandlerFunc func(c ContextWrap)

type Trace = trace.T

const (
	_Alias            = "_alias_"
	_TraceName        = "_trace_"
	_LoggerName       = "_logger_"
	_BodyName         = "_body_"
	_PayloadName      = "_payload_"
	_GraphPayloadName = "_graph_payload_"
	_SessionUserInfo  = "_session_user_info"
	_AbortErrorName   = "_abort_error_"
	_IsRecordMetrics  = "_is_record_metrics_"
)

type ContextWrap interface {
	init()

	// ShouldBindQuery 反序列化 querystring
	// tag: `form:"xxx"` (注：不要写成 query)
	ShouldBindQuery(obj interface{}) error

	// ShouldBindPostForm 反序列化 postform (querystring会被忽略)
	// tag: `form:"xxx"`
	ShouldBindPostForm(obj interface{}) error

	// ShouldBindForm 同时反序列化 querystring 和 postform;
	// 当 querystring 和 postform 存在相同字段时，postform 优先使用。
	// tag: `form:"xxx"`
	ShouldBindForm(obj interface{}) error

	// ShouldBindJSON 反序列化 postjson
	// tag: `json:"xxx"`
	ShouldBindJSON(obj interface{}) error

	// ShouldBindURI 反序列化 path 参数(如路由路径为 /user/:name)
	// tag: `uri:"xxx"`
	ShouldBindURI(obj interface{}) error

	// Redirect 重定向
	Redirect(code int, location string)

	// Trace 获取 Trace 对象
	Trace() Trace
	setTrace(trace Trace)
	disableTrace()

	// Logger 获取 Logger 对象
	Logger() *zap.Logger
	setLogger(logger *zap.Logger)

	// Payload 正确返回
	Payload(payload interface{})
	getPayload() interface{}

	// GraphPayload GraphQL返回值 与 api 返回结构不同
	GraphPayload(payload interface{})
	getGraphPayload() interface{}

	// HTML 返回界面
	HTML(name string, obj interface{})

	// AbortWithError 错误返回
	AbortWithError(err BusinessError)
	abortError() BusinessError

	// Header 获取 Header 对象
	Header() http.Header
	// GetHeader 获取 Header
	GetHeader(key string) string
	// SetHeader 设置 Header
	SetHeader(key, value string)

	// SessionUserInfo 当前用户信息
	SessionUserInfo() proposal.SessionUserInfo
	setSessionUserInfo(info proposal.SessionUserInfo)

	// Alias 设置路由别名 for metrics path
	Alias() string
	setAlias(path string)

	// disableRecordMetrics 设置禁止记录指标
	disableRecordMetrics()
	ableRecordMetrics()
	isRecordMetrics() bool

	// RequestInputParams 获取所有参数
	RequestInputParams() url.Values
	// RequestPostFormParams  获取 PostForm 参数
	RequestPostFormParams() url.Values
	// Request 获取 Request 对象
	Request() *http.Request
	// RawData 获取 Request.Body
	RawData() []byte
	// Method 获取 Request.Method
	Method() string
	// Host 获取 Request.Host
	Host() string
	// Path 获取 请求的路径 Request.URL.Path (不附带 querystring)
	Path() string
	// URI 获取 unescape 后的 Request.URL.RequestURI()
	URI() string
	// RequestContext 获取请求的 context (当 client 关闭后，会自动 canceled)
	RequestContext() StdContext

	// ResponseWriter 获取 ResponseWriter 对象
	ResponseWriter() gin.ResponseWriter
}

type GinContext struct {
	ctx *gin.Context
}

type StdContext struct {
	innerctx.Context
	Trace
	*zap.Logger
}

func (c *GinContext) init() {
	body, err := c.ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	c.ctx.Set(_BodyName, body)                                   // cache body是为了trace使用
	c.ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // re-construct req body
}

// ShouldBindQuery 反序列化querystring
// tag: `form:"xxx"` (注：不要写成query)
func (c *GinContext) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Query)
}

// ShouldBindPostForm 反序列化 postform (querystring 会被忽略)
// tag: `form:"xxx"`
func (c *GinContext) ShouldBindPostForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.FormPost)
}

// ShouldBindForm 同时反序列化querystring和postform;
// 当querystring和postform存在相同字段时，postform优先使用。
// tag: `form:"xxx"`
func (c *GinContext) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}

// ShouldBindJSON 反序列化postjson
// tag: `json:"xxx"`
func (c *GinContext) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}

// ShouldBindURI 反序列化path参数(如路由路径为 /user/:name)
// tag: `uri:"xxx"`
func (c *GinContext) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

// Redirect 重定向
func (c *GinContext) Redirect(code int, location string) {
	c.ctx.Redirect(code, location)
}

func (c *GinContext) Trace() Trace {
	t, ok := c.ctx.Get(_TraceName)
	if !ok || t == nil {
		return nil
	}

	return t.(Trace)
}

func (c *GinContext) setTrace(trace Trace) {
	c.ctx.Set(_TraceName, trace)
}

func (c *GinContext) disableTrace() {
	c.setTrace(nil)
}

func (c *GinContext) Logger() *zap.Logger {
	logger, ok := c.ctx.Get(_LoggerName)
	if !ok {
		return nil
	}

	return logger.(*zap.Logger)
}

func (c *GinContext) setLogger(logger *zap.Logger) {
	c.ctx.Set(_LoggerName, logger)
}

func (c *GinContext) getPayload() interface{} {
	if payload, ok := c.ctx.Get(_PayloadName); ok != false {
		return payload
	}
	return nil
}

func (c *GinContext) Payload(payload interface{}) {
	c.ctx.Set(_PayloadName, payload)
}

func (c *GinContext) getGraphPayload() interface{} {
	if payload, ok := c.ctx.Get(_GraphPayloadName); ok != false {
		return payload
	}
	return nil
}

func (c *GinContext) GraphPayload(payload interface{}) {
	c.ctx.Set(_GraphPayloadName, payload)
}

func (c *GinContext) HTML(name string, obj interface{}) {
	c.ctx.HTML(http.StatusOK, name+".html", obj)
}

func (c *GinContext) Header() http.Header {
	header := c.ctx.Request.Header

	clone := make(http.Header, len(header))
	for k, v := range header {
		value := make([]string, len(v))
		copy(value, v)

		clone[k] = value
	}
	return clone
}

func (c *GinContext) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *GinContext) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

func (c *GinContext) SessionUserInfo() proposal.SessionUserInfo {
	val, ok := c.ctx.Get(_SessionUserInfo)
	if !ok {
		return proposal.SessionUserInfo{}
	}

	return val.(proposal.SessionUserInfo)
}

func (c *GinContext) setSessionUserInfo(info proposal.SessionUserInfo) {
	c.ctx.Set(_SessionUserInfo, info)
}

func (c *GinContext) AbortWithError(err BusinessError) {
	if err != nil {
		httpCode := err.HTTPCode()
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}

		c.ctx.AbortWithStatus(httpCode)
		c.ctx.Set(_AbortErrorName, err)
	}
}

func (c *GinContext) abortError() BusinessError {
	err, _ := c.ctx.Get(_AbortErrorName)
	return err.(BusinessError)
}

func (c *GinContext) Alias() string {
	path, ok := c.ctx.Get(_Alias)
	if !ok {
		return ""
	}

	return path.(string)
}

func (c *GinContext) setAlias(path string) {
	if path = strings.TrimSpace(path); path != "" {
		c.ctx.Set(_Alias, path)
	}
}

func (c *GinContext) isRecordMetrics() bool {
	isRecordMetrics, ok := c.ctx.Get(_IsRecordMetrics)
	if !ok {
		return false
	}

	return isRecordMetrics.(bool)
}

func (c *GinContext) ableRecordMetrics() {
	c.ctx.Set(_IsRecordMetrics, true)
}

func (c *GinContext) disableRecordMetrics() {
	c.ctx.Set(_IsRecordMetrics, false)
}

// RequestInputParams 获取所有参数
func (c *GinContext) RequestInputParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.Form
}

// RequestPostFormParams 获取 PostForm 参数
func (c *GinContext) RequestPostFormParams() url.Values {
	_ = c.ctx.Request.ParseForm()
	return c.ctx.Request.PostForm
}

// Request 获取 Request
func (c *GinContext) Request() *http.Request {
	return c.ctx.Request
}

func (c *GinContext) RawData() []byte {
	body, ok := c.ctx.Get(_BodyName)
	if !ok {
		return nil
	}

	return body.([]byte)
}

// Method 请求的method
func (c *GinContext) Method() string {
	return c.ctx.Request.Method
}

// Host 请求的host
func (c *GinContext) Host() string {
	return c.ctx.Request.Host
}

// Path 请求的路径(不附带querystring)
func (c *GinContext) Path() string {
	return c.ctx.Request.URL.Path
}

// URI unescape后的uri
func (c *GinContext) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

// RequestContext (包装 Trace + Logger) 获取请求的 context (当client关闭后，会自动canceled)
func (c *GinContext) RequestContext() StdContext {
	return StdContext{
		//c.ctx.Request.Context(),
		innerctx.Background(),
		c.Trace(),
		c.Logger(),
	}
}

// ResponseWriter 获取 ResponseWriter
func (c *GinContext) ResponseWriter() gin.ResponseWriter {
	return c.ctx.Writer
}

var ginContextPool = &sync.Pool{
	New: func() interface{} {
		return new(GinContext)
	},
}

func newContext(ctx *gin.Context) ContextWrap {
	context := ginContextPool.Get().(*GinContext)
	context.ctx = ctx
	return context
}

func releaseContext(ctx ContextWrap) {
	c := ctx.(*GinContext)
	c.ctx = nil
	ginContextPool.Put(c)
}
