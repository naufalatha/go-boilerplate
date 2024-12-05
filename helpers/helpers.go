package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/naufalatha/go-boilerplate/helpers/logger"
	"github.com/naufalatha/go-boilerplate/models"
	"golang.org/x/exp/rand"
)

// GetClaim function to retrieve claims from JWT token in Fiber context
func GetClaim(c *fiber.Ctx) jwt.MapClaims {
	user := c.Locals("user")
	claims := jwt.MapClaims{}
	if user != nil {
		token := user.(*jwt.Token)
		claims = token.Claims.(jwt.MapClaims)
	}
	return claims
}

// GetCustomerIDFromClaim returns customer id from subject claim
func GetCustomerIDFromClaim(c *fiber.Ctx) (int64, error) {
	claims := GetClaim(c)
	if claims["sub"] != nil {
		var customerID int64
		var err error

		switch claims["sub"].(type) {
		case string:
			customerID, err = strconv.ParseInt(fmt.Sprintf("%v", claims["sub"]), 0, 64)
			if err != nil {
				logger.ErrorCaller(err.Error(), 0)
				return 0, fmt.Errorf("invalid access token")
			}
		case float64:
			sub, err := strconv.ParseFloat(fmt.Sprintf("%v", claims["sub"]), 64)
			if err != nil {
				logger.ErrorCaller(err.Error(), 0)
				return 0, fmt.Errorf("invalid access token")
			}
			customerID = int64(sub)
		default:
			logger.Error("unhandled customer id type")
			return 0, fmt.Errorf("invalid access token")
		}

		return customerID, nil
	}
	return 0, nil
}

// Generate new claim by customer id
func GenerateNewTokenClaims(customerID int64, expiry time.Duration, algorithm string) *jwt.Token {
	claims := jwt.New(jwt.GetSigningMethod(algorithm))
	claims.Claims = jwt.RegisteredClaims{
		ID:        uuid.New().String(),
		Subject:   fmt.Sprint(customerID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	return claims
}

// GetJSONTag extract json tag name from struct
func GetJSONTag(field reflect.StructField) string {
	splitN := strings.SplitN(field.Tag.Get("json"), ",", 2)
	if len(splitN) > 0 {
		name := splitN[0]
		if name == "-" {
			return ""
		}
		return name
	}
	return field.Tag.Get("json")
}

func ParseStringSnakeCase(value string) string {
	converted := regexp.MustCompile(`\"([^\"]*?)\"\s*?:`).ReplaceAllString(value, "${1}_${2}")
	return strings.ToLower(regexp.MustCompile(`(\w[^A-Z])([A-Z])`).ReplaceAllString(converted, "${1}_${2}"))
}

// CleanPhoneNumber transform phone number raw input into alphanumeric only
func CleanPhoneNumber(phoneNumber string) string {
	// return regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(phoneNumber, "")
	phone := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(phoneNumber, "")
	// replace if first character is 0
	if phone[0] == '0' {
		phone = "62" + phone[1:]
	} else if phone[0] != '6' || phone[1] != '2' {
		// If the number does not start with '62', add '62' at the beginning
		phone = "62" + phone
	}

	return phone
}

// FiberError simplify fiber error function
func FiberError(err string, code ...int) *fiber.Error {
	logger.ErrorCaller(err, 2)
	logger.Error(err)
	switch true {
	case err == models.ErrNoRows.Error() || strings.Contains(strings.ToLower(err), "not found") || strings.Contains(strings.ToLower(err), "tidak ditemukan"):
		return fiber.NewError(404, "Data tidak ditemukan")
	case strings.Contains(err, "pq:") || strings.Contains(err, "strconv"):
		return fiber.NewError(500, "Terjadi kesalahan pada sistem")
	case strings.Contains(strings.ToLower(err), "timeout") || strings.Contains(strings.ToLower(err), "host"):
		return fiber.NewError(500, "Sedang terjadi gangguan jaringan. Mohon coba kembali beberapa saat lagi.")
	default:
		statusCode := fiber.StatusInternalServerError
		if len(code) > 0 {
			statusCode = code[0]
		}
		return fiber.NewError(statusCode, "Sedang ada gangguan di jaringan atau sistem. Mohon coba kembali beberapa saat lagi.")
	}
}

// FiberErrorCustom simplify fiber error function
func FiberErrorCustom(message string, code ...int) *fiber.Error {
	logger.ErrorCaller(message, 2)
	logger.Error(message)
	statusCode := fiber.StatusInternalServerError
	if len(code) > 0 {
		statusCode = code[0]
	}
	return fiber.NewError(statusCode, message)
}

// FiberSuccess simplify fiber success function
func FiberSuccess(c *fiber.Ctx, data interface{}, messages ...string) error {
	message := "Data berhasil didapatkan"
	if len(messages) > 0 {
		message = messages[0]
	}
	err := c.JSON(models.Response{
		Success:    true,
		StatusCode: fiber.StatusOK,
		Data:       data,
		Message:    message,
	})
	logger.TraceCtx(c, "request returned successfully")

	return err
}

// FiberSuccessWithMessage simplify fiber success function with custom message without data
func FiberSuccessWithMessage(c *fiber.Ctx, message string) error {
	return c.JSON(models.Response{
		Success:    true,
		StatusCode: fiber.StatusOK,
		Message:    message,
	})
}

func FiberSuccessWithStatusCode(c *fiber.Ctx, message string, statusCode int) error {
	return c.JSON(models.Response{
		Success:    true,
		StatusCode: statusCode,
		Message:    message,
	})
}

// FiberFailure simplify fiber empty data function
func FiberEmptyData(c *fiber.Ctx) error {
	return c.JSON(models.Response{
		Success:    true,
		StatusCode: fiber.StatusOK,
		Message:    "Data tidak ditemukan!",
		Data: models.Pagination{
			Data: []string{},
		},
	})
}

func FiberErrorWithEmptyData(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Data tidak ditemukan!"
	}

	return c.JSON(models.Response{
		Status:     "ERROR",
		StatusCode: fiber.StatusInternalServerError,
		Message:    message,
		Data: models.Pagination{
			Data: []string{},
		},
	})
}

// ParsePaginationFilter help to parse pagination into format
func ParsePaginationFilter(queries map[string]string) models.PaginationFilter {
	pageParam := queries["page"]
	page, _ := strconv.ParseInt(pageParam, 0, 64)
	if page == 1 {
		page = 0
	}
	if page > 0 {
		page = (page - 1)
	}

	limitParam := queries["limit"]
	limit, _ := strconv.ParseInt(limitParam, 0, 64)
	if limit == 0 {
		limit = 10
	}

	return models.PaginationFilter{
		Limit: limit,
		Page:  page * limit,
	}
}

// ParseQueryParams help to parse request query params into mapped struct
func ParseQueryParams(c *fiber.Ctx, destination interface{}) error {
	if reflect.ValueOf(destination).Kind() != reflect.Ptr {
		return fmt.Errorf("ParseQueryParams: destination must be a pointer")
	}

	queries := c.Queries()

	t := reflect.TypeOf(destination).Elem()
	v := reflect.ValueOf(destination).Elem()

	for i := 0; i < v.NumField(); i++ {
		key := ParseStringSnakeCase(t.Field(i).Name)
		if key == "pagination_filter" {
			pagination := ParsePaginationFilter(queries)
			ReflectFormatValue(v.Field(i), pagination)
		}
		if key == "sort_by" {
			if len(queries["sortby"]) != 0 {
				ReflectFormatValue(v.Field(i), queries["sortby"])
			}
			if len(queries["sort_by"]) != 0 {
				ReflectFormatValue(v.Field(i), queries["sort_by"])
			}
		}

		if key == "keyword" {
			keyword := queries["keyword"]
			if len(keyword) > 3 {
				ReflectFormatValue(v.Field(i), queries["keyword"])
			}
		}

		if key == "customer_id" {
			customerID, err := GetCustomerIDFromClaim(c)
			if err != nil {
				logger.Error(err.Error())
				return err
			}
			ReflectFormatValue(v.Field(i), customerID)
		}

		key = fmt.Sprintf("filter[%s]", ParseStringSnakeCase(t.Field(i).Name))
		if value, ok := queries[key]; ok {
			ReflectFormatValue(v.Field(i), value)
		}
	}
	return nil
}

// ReflectFormatValue Set reflect field value according to registered types
// Use this when you want to set value automatically using the source data types
func ReflectFormatValue(field reflect.Value, value interface{}) error {
	if field.IsValid() && field.CanSet() {
		switch field.Kind() {
		case reflect.String:
			field.SetString(fmt.Sprint(value))
		case reflect.Int, reflect.Int64:
			if val, err := strconv.ParseInt(fmt.Sprint(value), 0, 64); err == nil {
				field.SetInt(val)
			}
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(fmt.Sprint(value), 64)
			if err != nil {
				return err
			}
			field.SetFloat(val)
		case reflect.Struct: // non-default type
			if field.Type() == reflect.TypeOf(models.PaginationFilter{}) {
				field.Set(reflect.ValueOf(value))
			}
		default:
			return fmt.Errorf("unsupported types")
		}
	}
	return nil
}

// HashSHA256 hash given input combined with salt
func HashSHA256(content []byte) string {
	hash := sha256.New()
	hash.Write(content)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// HashHMAC hash hmac using given input
func HashHMAC(message, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(message)
	return hex.EncodeToString(h.Sum(nil))
}

// Hash Security Code (PIN) using SHA256 and phone number
func HashSecurityCode(phone, securityCode string) string {
	return HashSHA256([]byte(phone + securityCode))
}

// IsNumeric check if given string is numeric
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// checking if string contains special character (+, -, /, \, *, etc)
func ContainsOtherCharacter(s string) bool {
	return strings.ContainsAny(s, "~!@#$%^&*()_+`-={}|[]\\:\";'<>?,./")
}

// ContainsSpecialCharacter check if given string contains special character
func ContainsSpecialCharacter(s string) bool {
	return strings.ContainsAny(s, "\r\n")
}

func ModifyAddress(address string) string {
	// Pisahkan alamat berdasarkan koma
	parts := strings.Split(address, ", ")
	if len(parts) > 0 && strings.Contains(parts[0], "+") && len(parts[0]) <= 10 {
		// Hapus bagian pertama dari alamat
		address = strings.Join(parts[1:], ", ")
		// Jika karakter pertama adalah spasi, hapus spasi tersebut
		if len(address) > 0 && address[0] == ' ' {
			address = strings.TrimPrefix(address, " ")
		}
	}
	return address
}

func GenerateTokenResetPassword(length int) string {
	characters := "0123456789"
	charactersLength := len(characters)
	randomString := ""

	for i := 0; i < length; i++ {
		randomString += string(characters[rand.Intn(charactersLength)])
	}

	return randomString
}

func GenerateUsernameFromEmail(email string) string {
	username := strings.Split(email, "@")[0]
	randomNumber := strconv.Itoa(rand.Intn(999-100) + 100)
	return username + randomNumber
}
