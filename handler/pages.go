package handler

import (
	"log"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mccune1224/betrayal-widget/data"
	"github.com/mccune1224/betrayal-widget/util"
	"github.com/mccune1224/betrayal-widget/view/page"
)

func (h *Handler) IndexPage(c echo.Context) error {
	games, err := h.models.Games.GetAll()
	if err != nil {
		log.Println(err)
		return TemplRender(c, page.Error500(c, err))
	}

	return TemplRender(c, page.Index(c, games))
}

func (h *Handler) GameDashboardPage(c echo.Context) error {
	gID, err := strconv.Atoi(c.Param("game_id"))
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	game, err := h.models.Games.GetByID(gID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	return TemplRender(c, page.GameDashboard(c, game))
}

func (h *Handler) JoinGamePage(c echo.Context) error {
	gID, err := strconv.Atoi(c.Param("game_id"))
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	game, err := h.models.Games.GetByID(gID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	players, err := h.models.Players.GetAllByGameID(game.GameID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	return TemplRender(c, page.JoinGame(c, game, players))
}

func (h *Handler) PlayerDashboardPage(c echo.Context) error {
	game, _ := util.GetGame(c)
	playerID := util.ParamInt(c, "player_id", -1)

	player, err := h.models.Players.GetByID(playerID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}
	role, err := h.models.Roles.GetComplexByID(player.RoleID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	actions, err := h.models.Actions.GetAll()
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	for _, roleAbility := range role.Abilities {
		actions = append(actions, data.Action{
			AbilityName:     roleAbility.Name + " (provided by role)",
			Description:     roleAbility.Description,
			Rarity:          roleAbility.Rarity,
			RoleAssociation: role.Name,
			Categories:      roleAbility.Categories,
		})
	}

	sort.Slice(actions, func(i, j int) bool {
		return actions[i].AbilityName < actions[j].AbilityName
	})

	pa, err := h.models.Actions.GetAllPlayerUnapprovedRequestsByPlayerID(player.ID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	return TemplRender(c, page.PlayerDashboard(c, game, player, role, actions, pa))
}

func (h *Handler) PlayerFlashcard(c echo.Context) error {
	playerID := util.ParamInt(c, "player_id", -1)

	player, err := h.models.Players.GetByID(playerID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}
	playerRole, err := h.models.Roles.GetComplexByID(player.RoleID)
	if err != nil {
		return TemplRender(c, page.Error500(c, err))
	}

	return TemplRender(c, page.PlayerFlashcard(c, player, playerRole, nil))
}
