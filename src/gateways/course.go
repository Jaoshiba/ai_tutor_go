package gateways

import (
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) CreateCourse(ctx *fiber.Ctx) error {

	jsonbody := ctx.FormValue("jsonbody")
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid file"})
	}

	fmt.Println("jsonbody: ", jsonbody)

	var coursejsonBody entities.CourseRequestBody

	err = json.Unmarshal([]byte(jsonbody), &coursejsonBody)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	if coursejsonBody.Title == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed course name"})
	} else if coursejsonBody.Description == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed description"})
	}

	fmt.Println("Before create in gateway")

	courses, err := h.CourseService.CreateCourse(coursejsonBody, false, file, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "failed to create course on CreateCourse",
			Data:    err.Error(),
		})
	}

	datareturn := map[string]interface{}{
		"confirmed":   coursejsonBody.Confirmed,
		"isfirsttime": coursejsonBody.IsFirtTime,
		"courses":     courses,
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "Completed create Course from your promts",
		Data:    datareturn,
	})
}

func (h *HTTPGateway) GetCourseByUser(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "User ID not found in context"})
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "Invalid user ID in context"})
	}

	courses, err := h.CourseService.GetCourses(c)

	fmt.Println(courses)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "Failed to retrieve courses"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Courses retrieved successfully",
		"data":    courses,
	})
}

func (h *HTTPGateway) GetCourseDetail(c *fiber.Ctx) error {
	fmt.Println("Hello im in gateway at start")
	// 1. ดึง CourseID จาก URL parameter
	// สมมติว่า URL endpoint คือ /api/v1/courses/:courseId
	courseID := c.Params("courseId")
	if courseID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Course ID is required")
	}

	fmt.Printf("Attempting to get details for Course ID: %s\n", courseID)

	fmt.Println("Hello im in gateway at start")
	// 2. เรียกใช้ CourseService.GetCourseDetail
	courseDetail, err := h.CourseService.GetCourseDetail(c, courseID)
	if err != nil {
		fmt.Println("Hello im in gateway at start")
		// จัดการข้อผิดพลาดจาก Service Layer
		// Fiber.Error จะถูกแปลงเป็น HTTP status code โดย Fiber middleware อัตโนมัติ
		if fiberErr, ok := err.(*fiber.Error); ok {
			fmt.Printf("Service error for Course ID %s: %s (Status: %d)\n", courseID, fiberErr.Message, fiberErr.Code)
			return fiberErr
		}
		// สำหรับ error ทั่วไปที่ไม่ใช่ Fiber.Error ให้ return Internal Server Error
		fmt.Printf("Unexpected error getting course details for Course ID %s: %v\n", courseID, err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to get course details: %v", err))
	}

	// 3. ส่งคืนข้อมูลในรูปแบบ JSON
	fmt.Printf("Successfully retrieved details for Course ID: %s\n", courseID)
	return c.Status(fiber.StatusOK).JSON(courseDetail)
}

func (h *HTTPGateway) DeleteCourse(c *fiber.Ctx) error {

	courseID := c.Params("courseId")
	if courseID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Course ID is required")
	}

	err := h.CourseService.DeleteCourse(c, courseID)
	if err != nil {
		fmt.Println("Error deleting course:", err) // แนะนำให้ log ข้อผิดพลาด
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete course")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Course deleted successfully",
	})
}
