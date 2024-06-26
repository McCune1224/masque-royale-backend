package handler

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/mccune1224/betrayal-widget/models"
	"github.com/mccune1224/betrayal-widget/util"
)

func (h *Handler) GetAllRoles(c echo.Context) error {

	q := models.New(h.Db)

	roles, err := q.ListRoles(c.Request().Context())
	if err != nil {
		return c.JSON(500, err)
	}

	return c.JSON(200, roles)
}

func (h *Handler) GetRoleByID(c echo.Context) error {
	roleID, err := util.ParseInt32Param(c, "role_id")
	if err != nil {
		return util.BadRequestJson(c, "Invalid Role ID")
	}
	q := models.New(h.Db)
	role, err := q.GetRole(c.Request().Context(), roleID)
	if err != nil {
		log.Println(err)
		return util.InternalServerErrorJson(c, err.Error())
	}

	return c.JSON(200, role)
}

type CompleteRole struct {
	Role      models.Role            `json:"role"`
	Passives  []models.PassiveDetail `json:"passives"`
	Abilities []models.AbilityDetail `json:"abilities"`
}

func (h *Handler) GetCompleteRole(c echo.Context) error {
	roleID, err := util.ParseInt32Param(c, "role_id")
	if err != nil {
		return util.BadRequestJson(c, "Invalid Role ID")
	}
	q := models.New(h.Db)
	role, err := q.GetRole(c.Request().Context(), roleID)
	if err != nil {
		log.Println(err)
		return util.InternalServerErrorJson(c, err.Error())
	}

	passives, err := q.GetAssociatedRolePassives(context.Background(), roleID)
	if err != nil {
		return util.InternalServerErrorJson(c, err.Error())
	}

	abilities, err := q.GetRoleAbilityDetails(context.Background(), roleID)
	if err != nil {
		log.Println(err)
		return util.InternalServerErrorJson(c, err.Error())
	}
	completeRole := CompleteRole{Role: role, Passives: passives, Abilities: abilities}
	return c.JSON(200, completeRole)
}

func (h *Handler) GetAllCompleteRoles(c echo.Context) error {
	q := models.New(h.Db)
	roles, err := q.ListRoles(c.Request().Context())
	if err != nil {
		return util.InternalServerErrorJson(c, err.Error())
	}
	completeRoles := []CompleteRole{}
	for _, role := range roles {
		passives, err := q.GetAssociatedRolePassives(context.Background(), role.ID)
		if err != nil {
			return util.InternalServerErrorJson(c, err.Error())
		}
		abilities, err := q.GetRoleAbilityDetails(context.Background(), role.ID)
		if err != nil {
			return util.InternalServerErrorJson(c, err.Error())
		}
		completeRoles = append(completeRoles, CompleteRole{Role: role, Passives: passives, Abilities: abilities})
	}
	return c.JSON(200, completeRoles)
}
func (h *Handler) GetRoleAbilities(c echo.Context) error {
	roleId, err := util.ParseInt32Param(c, "role_id")
	if err != nil {
		return util.BadRequestJson(c, "Invalid Role ID")
	}
	q := models.New(h.Db)
	abilities, err := q.GetRoleAbilityDetails(c.Request().Context(), roleId)
	if err != nil {
		return util.InternalServerErrorJson(c, err.Error())
	}
	return c.JSON(200, abilities)
}

func (h *Handler) GetRolePassives(c echo.Context) error {
	roleId, err := util.ParseInt32Param(c, "role_id")
	if err != nil {
		return util.BadRequestJson(c, "Invalid Role ID")
	}
	q := models.New(h.Db)
	passives, err := q.GetRolePassiveDetails(c.Request().Context(), roleId)
	if err != nil {
		return util.InternalServerErrorJson(c, err.Error())
	}
	return c.JSON(200, passives)
}
