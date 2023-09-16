package rest

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vamika-digital/wms-api-server/config/registry"
	"github.com/vamika-digital/wms-api-server/global/drivers"
	"github.com/vamika-digital/wms-api-server/interface/rest/admin"
	"github.com/vamika-digital/wms-api-server/interface/rest/auth"
	"github.com/vamika-digital/wms-api-server/interface/rest/master"
	"github.com/vamika-digital/wms-api-server/interface/rest/mobile"
	"github.com/vamika-digital/wms-api-server/interface/rest/warehouse"

	swaggerFiles "github.com/swaggo/files"
)

// swagger embed files

type Server struct {
	Address         string
	Port            int
	AdminModule     *admin.AdminRestModule
	AuthModule      *auth.AuthRestModule
	MasterModule    *master.MasterRestModule
	WarehouseModule *warehouse.WarehouseRestModule
	MobileModule    *mobile.MobileRestModule
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
	joborderService := registry.GetContainerInstance(conn).JobOrderService
	outwardrequestService := registry.GetContainerInstance(conn).OutwardRequestService
	requisitionService := registry.GetContainerInstance(conn).RequisitionService

	batchlabelService := registry.GetContainerInstance(conn).BatchLabelService
	labelstickerService := registry.GetContainerInstance(conn).LabelStickerService
	inventoryService := registry.GetContainerInstance(conn).InventoryService
	stockService := registry.GetContainerInstance(conn).StockService

	authModule := auth.NewAuthRestModule(authService)
	adminModule := admin.NewAdminRestModule(roleService, userService)
	masterModule := master.NewMasterRestModule(machineService, customerService, storeService, containerService, productService, joborderService, outwardrequestService, requisitionService)
	warehouseModule := warehouse.NewWarehouseRestModule(batchlabelService, labelstickerService, inventoryService, stockService)
	mobileRestModule := mobile.NewMobileRestModule(containerService, inventoryService, requisitionService, outwardrequestService)

	return &Server{
		Address:         address,
		Port:            port,
		AdminModule:     adminModule,
		AuthModule:      authModule,
		MasterModule:    masterModule,
		WarehouseModule: warehouseModule,
		MobileModule:    mobileRestModule,
	}
}

func (s *Server) Run() {
	r := setupRouter()
	rootGroup := r.Group("/")

	s.AdminModule.RegisterRoutes(rootGroup)
	s.AuthModule.RegisterRoutes(rootGroup)
	s.MasterModule.RegisterRoutes(rootGroup)
	s.WarehouseModule.RegisterRoutes(rootGroup)
	s.MobileModule.RegisterRoutes(rootGroup)

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
