package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"main/src/models"
	"main/src/repositories"
)

type TaskHandler struct {
	repo *repositories.TaskRepository
}

func NewTaskHandler(r *repositories.TaskRepository) *TaskHandler {
	return &TaskHandler{
		repo: r,
	}
}

// Вспомогательная функция для возврата ошибки.
func fail(c fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// Получение и проверка ID задачи.
func parseID(c fiber.Ctx) (int64, bool) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)

	if err != nil || id < 1 {
		fail(c, http.StatusBadRequest, "id must be a positive integer")
		return 0, false
	}

	return id, true
}

// Проверка допустимого статуса задачи.
func validStatus(s string) bool {
	return s == "new" ||
		s == "in_progress" ||
		s == "done"
}

// POST /api/tasks

// Create godoc
// @Summary Создать задачу
// @Description Создаёт новую задачу с указанным названием и описанием. Статус устанавливается по умолчанию.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Данные новой задачи"
// @Success 201 {object} models.Task "Задача успешно создана"
// @Failure 400 {object} map[string]string "Некорректный JSON или название задачи"
// @Failure 500 {object} map[string]string "Ошибка создания задачи"
// @Router /api/tasks [post]
func (h *TaskHandler) Create(c fiber.Ctx) error {
	var req models.CreateTaskRequest

	if err := c.Bind().Body(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid JSON")
	}

	req.Title = strings.TrimSpace(req.Title)

	if req.Title == "" || len(req.Title) > 120 {
		return fail(
			c,
			http.StatusBadRequest,
			"title is required and must be at most 120 characters",
		)
	}

	t, err := h.repo.Create(
		c,
		req.Title,
		req.Description,
	)

	if err != nil {
		return fail(
			c,
			http.StatusInternalServerError,
			"could not create task",
		)
	}

	return c.Status(http.StatusCreated).JSON(t)
}

// GET /api/tasks

// List godoc
// @Summary Получить список задач
// @Description Возвращает список всех задач.
// @Tags Tasks
// @Produce json
// @Success 200 {array} models.Task "Список задач"
// @Failure 500 {object} map[string]string "Ошибка получения списка задач"
// @Router /api/tasks [get]
func (h *TaskHandler) List(c fiber.Ctx) error {
	tasks, err := h.repo.List(c)

	if err != nil {
		return fail(
			c,
			http.StatusInternalServerError,
			"could not retrieve tasks",
		)
	}

	return c.Status(http.StatusOK).JSON(tasks)
}

// GET /api/tasks/:id

// Get godoc
// @Summary Получить задачу
// @Description Возвращает задачу по её идентификатору.
// @Tags Tasks
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Success 200 {object} models.Task "Задача найдена"
// @Failure 400 {object} map[string]string "Некорректный идентификатор"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Ошибка получения задачи"
// @Router /api/tasks/{id} [get]
func (h *TaskHandler) Get(c fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return nil
	}

	t, err := h.repo.Get(c, id)

	if errors.Is(err, repositories.ErrNotFound) {
		return fail(c, http.StatusNotFound, "task not found")
	}

	if err != nil {
		return fail(
			c,
			http.StatusInternalServerError,
			"could not retrieve task",
		)
	}

	return c.Status(http.StatusOK).JSON(t)
}

// PUT /api/tasks/:id

// Update godoc
// @Summary Обновить задачу
// @Description Полностью обновляет название, описание и статус существующей задачи.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Param task body models.UpdateTaskRequest true "Новые данные задачи"
// @Success 200 {object} models.Task "Задача успешно обновлена"
// @Failure 400 {object} map[string]string "Некорректный JSON, идентификатор или статус"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Ошибка обновления задачи"
// @Router /api/tasks/{id} [put]
func (h *TaskHandler) Update(c fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return nil
	}

	var req models.UpdateTaskRequest

	if err := c.Bind().Body(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid JSON")
	}

	req.Title = strings.TrimSpace(req.Title)

	if req.Title == "" || len(req.Title) > 120 {
		return fail(
			c,
			http.StatusBadRequest,
			"title is required and must be at most 120 characters",
		)
	}

	if !validStatus(req.Status) {
		return fail(
			c,
			http.StatusBadRequest,
			"status must be new, in_progress or done",
		)
	}

	t, err := h.repo.Update(
		c,
		id,
		req.Title,
		req.Description,
		req.Status,
	)

	if errors.Is(err, repositories.ErrNotFound) {
		return fail(c, http.StatusNotFound, "task not found")
	}

	if err != nil {
		return fail(
			c,
			http.StatusInternalServerError,
			"could not update task",
		)
	}

	return c.Status(http.StatusOK).JSON(t)
}

// PATCH /api/tasks/:id/status

// UpdateStatus godoc
// @Summary Изменить статус задачи
// @Description Изменяет только статус существующей задачи.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path int true "Идентификатор задачи"
// @Param status body models.UpdateStatusRequest true "Новый статус задачи"
// @Success 200 {object} models.Task "Статус успешно обновлён"
// @Failure 400 {object} map[string]string "Некорректный JSON, идентификатор или статус"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Ошибка обновления статуса"
// @Router /api/tasks/{id}/status [patch]
func (h *TaskHandler) UpdateStatus(c fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return nil
	}

	var req models.UpdateStatusRequest

	if err := c.Bind().Body(&req); err != nil {
		return fail(c, http.StatusBadRequest, "invalid JSON")
	}

	if !validStatus(req.Status) {
		return fail(
			c,
			http.StatusBadRequest,
			"status must be new, in_progress or done",
		)
	}

	t, err := h.repo.UpdateStatus(
		c,
		id,
		req.Status,
	)

	if errors.Is(err, repositories.ErrNotFound) {
		return fail(c, http.StatusNotFound, "task not found")
	}

	if err != nil {
		return fail(
			c,
			http.StatusInternalServerError,
			"could not update task status",
		)
	}

	return c.Status(http.StatusOK).JSON(t)
}

// DELETE /api/tasks/:id

// Delete godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по её идентификатору.
// @Tags Tasks
// @Param id path int true "Идентификатор задачи"
// @Success 204 "Задача успешно удалена"
// @Failure 400 {object} map[string]string "Некорректный идентификатор"
// @Failure 404 {object} map[string]string "Задача не найдена"
// @Failure 500 {object} map[string]string "Ошибка удаления задачи"
// @Router /api/tasks/{id} [delete]
func (h *TaskHandler) Delete(c fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return nil
	}

	err := h.repo.Delete(c, id)

	if errors.Is(err, repositories.ErrNotFound) {
		return fail(c, http.StatusNotFound, "task not found")
	}

	if err != nil {
		return fail(
			c,
			http.StatusInternalServerError,
			"could not delete task",
		)
	}

	return c.SendStatus(http.StatusNoContent)
}

