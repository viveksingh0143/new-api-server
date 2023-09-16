package stock

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	masterService "github.com/vamika-digital/wms-api-server/core/business/master/service"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/inventory"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
)

type StockRestHandler struct {
	InventoryService      service.InventoryService
	RequisitionService    masterService.RequisitionService
	OutwardRequestService masterService.OutwardRequestService
}

type ContainerAttachmentParams struct {
	SourceCode      string `json:"source_code" validation:"required"`
	DestinationCode string `json:"destination_code" validation:"required"`
}

func NewStockHandler(inventoryService service.InventoryService, requisitionService masterService.RequisitionService, outwardRequestService masterService.OutwardRequestService) *StockRestHandler {
	return &StockRestHandler{InventoryService: inventoryService, RequisitionService: requisitionService, OutwardRequestService: outwardRequestService}
}

func (h *StockRestHandler) CreateRawMaterialStock(c *gin.Context) {
	var formDTO = &inventory.InventoryRMStockCreateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	formDTO.CreatePallet = customtypes.NewValidNullBool(false)
	if err := h.InventoryService.CreateRawMaterialStock(formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *StockRestHandler) CreateFinishedStocks(c *gin.Context) {
	var formDTO = &inventory.InventoryFDStockCreateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	if err := h.InventoryService.CreateFinishedGoodsStock(formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *StockRestHandler) CreateFinishedStock(c *gin.Context) {
	var formDTO = &inventory.InventoryFDSingleStockCreateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	newFormDTO := &inventory.InventoryFDStockCreateDto{
		StoreID:       formDTO.StoreID,
		BinCode:       formDTO.BinCode,
		Barcodes:      []string{formDTO.Barcode},
		LastUpdatedBy: formDTO.LastUpdatedBy,
	}

	if err := h.InventoryService.CreateFinishedGoodsStock(newFormDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *StockRestHandler) GetAllContainerStocks(c *gin.Context) {
	var filter = &inventory.InventoryFilterDto{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	var pageProps dto.PaginationProps
	if err := c.ShouldBindQuery(&pageProps); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	pageNumber, rowsPerPage, sort := pageProps.GetValues()
	data, totalCount, err := h.InventoryService.GetAllProductsWithStockCounts(pageNumber, rowsPerPage, sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	pageResponse := dto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	c.JSON(http.StatusOK, pageResponse)
}

func (h *StockRestHandler) AttachContainer(c *gin.Context) {
	var containerParams = &ContainerAttachmentParams{}
	if err := c.ShouldBindJSON(&containerParams); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Invalid Inputs", nil))
		return
	}

	if err := validators.Validate.Struct(containerParams); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, containerParams)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	err := h.InventoryService.AttachContainer(containerParams.SourceCode, containerParams.DestinationCode)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, dto.GetRestResponse(http.StatusOK, "Successfully attached..."))
}

func (h *StockRestHandler) GetRequisitionByCode(c *gin.Context) {
	var order dto.OrderParams
	if err := c.ShouldBindQuery(&order); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Validate the struct
	if err := validators.Validate.Struct(order); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, order)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	requisition, reports, err := h.RequisitionService.GetRequisitionByCode(order.OrderNo)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, dto.GetRestDataResponse(http.StatusOK, "Requition with reports", map[string]interface{}{
		"requisition": requisition,
		"reports":     reports,
	}))
}
