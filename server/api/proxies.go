package api

import (
	"next-terminal/server/model"
	"next-terminal/server/utils"
	"strconv"

	"github.com/labstack/echo/v4"
)

// ProxiesAllEndpoint 返回全部的跳板代理
// 实际上最多只返回100个, 应该没人会设置那么多吧?
func ProxiesAllEndpoint(c echo.Context) error {
	items, total, err := proxyRepository.Find(1, 100, "", "", "", "")
	if err != nil {
		return err
	}

	return Success(c, H{
		"total": total,
		"items": items,
	})
}


// ProxiesGetEndpoint 返回给定ProxyID的跳板代理
func ProxiesGetEndpoint(c echo.Context) error {
	id := c.Param("id")

	item, err := proxyRepository.FindById(id)
	if err != nil {
		return err
	}

	return Success(c, item)
}


// ProxiesPagingEndpoint 基于给定的参数返回跳板代理列表
func ProxiesPagingEndpoint(c echo.Context) error {
	pageIndex, _ := strconv.Atoi(c.QueryParam("pageIndex"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	name := c.QueryParam("name")
	proxyType := c.QueryParam("proxyType")

	order := c.QueryParam("order")
	field := c.QueryParam("field")

	items, total, err := proxyRepository.Find(pageIndex, pageSize, name, proxyType, order, field)
	if err != nil {
		return err
	}

	return Success(c, H{
		"total": total,
		"items": items,
	})
}

// ProxiesCreateEndpoint 创建跳板代理
func ProxiesCreateEndpoint(c echo.Context) error {
	var item model.Proxy
	if err := c.Bind(&item); err != nil {
		return err
	}

	item.ID = utils.UUID()
	item.Created = utils.NowJsonTime()
	if !item.Type.Valid() {
		return Fail(c, -1, "暂未支持的Proxy类型")
	}

	if err := proxyRepository.Create(&item); err != nil {
		return err
	}

	return Success(c, "")
}

// ProxiesUpdateEndpoint 更新给定的Proxy
func ProxiesUpdateEndpoint(c echo.Context) error {
	id := c.Param("id")

	var item model.Proxy
	if err := c.Bind(&item); err != nil {
		return err
	}

	if err := proxyRepository.UpdateById(&item, id); err != nil {
		return err
	}

	return Success(c, nil)
}

// ProxiesDeleteEndpoint 基于给定的ProxyID删除Proxy
func ProxiesDeleteEndpoint(c echo.Context) error {
	id := c.Param("id")

	err := assetRepository.EmptyProxyByProxyID(id)
	if err != nil {
		return err
	}
	err = proxyRepository.DeleteByID(id)
	if err != nil {
		return err
	}

	return Success(c, nil)
}

// ProxiesUsageDetailAssetEndpoint 基于给定的ProxyID返回正在使用的资产列表
func ProxiesUsageDetailAssetEndpoint(c echo.Context) error {
	id := c.Param("id")

	assets, err := assetRepository.FindByProxyID(id)
	if err != nil {
		return err
	}

	return Success(c, H{
		"total": len(assets),
		"items": assets,
	})
}
