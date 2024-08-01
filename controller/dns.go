package controller

import (
	"github.com/cclose/dnsmasq-api/model"
	"github.com/cclose/dnsmasq-api/service"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type IDNSController interface {
	GetAllDNSRecords(ctx echo.Context) error
	GetDNSRecord(ctx echo.Context) error
	SetDNSRecord(ctx echo.Context) error
	DeleteDNSRecord(ctx echo.Context) error
	Register(e *echo.Echo)
}

type DnsController struct {
	ds service.IDNSMasqService
}

func NewDnsController(ds service.IDNSMasqService) IDNSController {
	return &DnsController{
		ds: ds,
	}
}

func (dc *DnsController) Register(e *echo.Echo) {
	e.GET("/dns", dc.GetAllDNSRecords)
	e.GET("/dns/:hostname", dc.GetDNSRecord)
	e.POST("/dns/:hostname", dc.SetDNSRecord)
	e.DELETE("/dns/:hostname", dc.DeleteDNSRecord)
}

func (dc *DnsController) GetAllDNSRecords(ctx echo.Context) error {
	records, err := dc.ds.GetAllIPs()
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	} // implicit else

	return ctx.JSON(http.StatusOK, records)
}

func (dc *DnsController) GetDNSRecord(ctx echo.Context) error {
	hostname := ctx.Param("hostname")
	records, err := dc.ds.GetIPByHost(hostname)
	if err != nil {
		if err.Error() == service.ErrorNoIPForHost {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": "hostname not found"})
		} // implicit else

		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	} // implicit else

	return ctx.JSON(http.StatusOK, records)
}

func (dc *DnsController) SetDNSRecord(ctx echo.Context) error {
	hostname := ctx.Param("hostname")
	// Parse the append query parameter
	appendStr := ctx.QueryParam("append")
	appendIP, err := strconv.ParseBool(appendStr)
	if err != nil {
		appendIP = false // Default to false if the parameter is not provided or invalid
	}

	req := model.SetDNSRecordRequest{}
	if err = ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body: " + err.Error()})
	}

	if len(req.IPs) == 0 {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "IP address list is required"})
	}

	records, err := dc.ds.SetIPByHost(hostname, req.IPs, appendIP)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err = dc.ds.UpdateDNSMasq(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, records)
}

func (dc *DnsController) DeleteDNSRecord(ctx echo.Context) error {
	hostname := ctx.Param("hostname")

	err := dc.ds.DeleteByHost(hostname)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if err = dc.ds.UpdateDNSMasq(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "hostname deleted"})
}
