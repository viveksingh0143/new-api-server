package rest

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/config/registry"
	"github.com/vamika-digital/wms-api-server/global/drivers"
	"github.com/vamika-digital/wms-api-server/interface/rest/admin"
	"github.com/vamika-digital/wms-api-server/interface/rest/auth"
	"github.com/vamika-digital/wms-api-server/interface/rest/master"
)

type Server struct {
	Address      string
	Port         int
	AdminModule  *admin.AdminRestModule
	AuthModule   *auth.AuthRestModule
	MasterModule *master.MasterRestModule
}

func NewServer(address string, port int, conn drivers.Connection) *Server {
	roleService := registry.GetContainerInstance(conn).RoleService
	userService := registry.GetContainerInstance(conn).UserService
	authService := registry.GetContainerInstance(conn).AuthService

	machineService := registry.GetContainerInstance(conn).MachineService
	customerService := registry.GetContainerInstance(conn).CustomerService
	storeService := registry.GetContainerInstance(conn).StoreService
	containerService := registry.GetContainerInstance(conn).ContainerService
	productService := registry.GetContainerInstance(conn).ProductService

	authModule := auth.NewAuthRestModule(authService)
	adminModule := admin.NewAdminRestModule(roleService, userService)
	masterModule := master.NewMasterRestModule(machineService, customerService, storeService, containerService, productService)

	return &Server{
		Address:      address,
		Port:         port,
		AdminModule:  adminModule,
		AuthModule:   authModule,
		MasterModule: masterModule,
	}
}

func (s *Server) Run() {
	r := setupRouter()
	rootGroup := r.Group("/")

	s.AdminModule.RegisterRoutes(rootGroup)
	s.AuthModule.RegisterRoutes(rootGroup)
	s.MasterModule.RegisterRoutes(rootGroup)

	log.Printf("Server started on %s:%d", s.Address, s.Port)
	r.Run(fmt.Sprintf("%s:%d", s.Address, s.Port))
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	return r
}
