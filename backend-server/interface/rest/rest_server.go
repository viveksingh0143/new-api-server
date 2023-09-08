package rest

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/config"
	"github.com/vamika-digital/wms-api-server/global/drivers"
	"github.com/vamika-digital/wms-api-server/interface/rest/admin"
)

type Server struct {
	Address     string
	Port        int
	AdminModule *admin.AdminRestModule
}

func NewServer(address string, port int, conn drivers.Connection) *Server {
	roleService := config.GetContainerInstance(conn).RoleService
	userService := config.GetContainerInstance(conn).UserService

	adminModule := admin.NewAdminRestModule(roleService, userService)

	// authModule := auth.NewAuthRestModule(db)
	// userModule := userRest.NewUserModule(db)
	// productModule := productRest.NewProductModule(db)
	// warehouseModule := warehouse.NewWarehouseModule(db)
	// return &Server{Address: address, Port: port, AuthModule: authModule, UserModule: userModule, ProductModule: productModule, WarehouseModule: warehouseModule}
	// return &Server{Address: address, Port: port, AuthModule: authModule}
	return &Server{Address: address, Port: port, AdminModule: adminModule}
}

func (s *Server) Run() {
	r := setupRouter()
	rootGroup := r.Group("/")

	s.AdminModule.RegisterRoutes(rootGroup)

	log.Printf("Server started on %s:%d", s.Address, s.Port)
	r.Run(fmt.Sprintf("%s:%d", s.Address, s.Port))
}

func setupRouter() *gin.Engine {
	// r := mux.NewRouter()

	// s.AuthModule.RegisterRoutes(r.PathPrefix("/").Subrouter())

	// r.Use(middlewares.ContentTypeMiddleware)
	// r.Use(middlewares.CORSMiddleware)
	// s.AuthModule.RegisterRoutes(r.PathPrefix("/auth").Subrouter())
	// s.UserModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	// s.ProductModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	// s.WarehouseModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())

	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	return r
}
