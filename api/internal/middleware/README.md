# middleware

```tree
middleware/
├── README.md
├── chain.go
│   ├── type Constructor func()
│   ├── type Chain {constructors: []Constructor}
│   ├── func New(constructors ...Constructor) Chain
│   ├── func (Chain) Then(h http.Handler) http.Handler
│   ├── func (Chain) ThenFunc(fn http.HandlerFunc) http.Handler
│   ├── func (Chain) Append(constructors ...Constructor) Chain
│   └── func (Chain) Extend(chain Chain) Chain
├── chain_test.go
│   ├── func tagMiddleware(tag string) Constructor
│   ├── func funcsEqual(f1 interface{}, f2 interface{}) bool
│   ├── func TestNew(t *testing.T)
│   ├── func TestThenWorksWithNoMiddleware(t *testing.T)
│   ├── func TestThenTreatsNilAsDefaultServeMux(t *testing.T)
│   ├── func TestThenFuncTreatsNilAsDefaultServeMux(t *testing.T)
│   ├── func TestThenFuncConstructsHandlerFunc(t *testing.T)
│   ├── func TestThenOrdersHandlersCorrectly(t *testing.T)
│   ├── func TestAppendAddsHandlersCorrectly(t *testing.T)
│   ├── func TestAppendRespectsImmutability(t *testing.T)
│   ├── func TestExtendAddsHandlersCorrectly(t *testing.T)
│   └── func TestExtendRespectsImmutability(t *testing.T)
└── cors.go
    ├── func allowedOrigins() []string
    ├── func isOriginAllowed(origin string, allowed []string) bool
    └── func CORS(next http.Handler) http.Handler
```
