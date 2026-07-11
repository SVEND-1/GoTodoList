package userApi

import (
	"TodoList/internal/core/domain"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/requests"
	"TodoList/internal/core/transport/http/response"
	core_types "TodoList/internal/core/transport/http/types"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type PatchUserRequest struct {
	FullName     core_types.Nullable[string] `json:"full_name"` //Валидировать можно только от go переменные,а не кастомные,но можно научить
	Phone_number core_types.Nullable[string] `json:"phone_number"`
}

func (r PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName can't be null")
		}
		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("FullName must be 3 and 100 characters")
		}
	}

	if r.Phone_number.Set {
		if r.Phone_number.Value != nil {
			phoneNumberLen := len([]rune(*r.Phone_number.Value))
			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf("PhoneNumber must be 10 or 15 characters")
			}
			if !strings.HasPrefix(*r.Phone_number.Value, "+") {
				return fmt.Errorf("PhoneNumber must start with '+'")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

func (c *UserController) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)
	responseHandler := response.NewHTTPResponseHandler(log, rw)

	userId, err := requests.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to parse userId",
		)
		return
	}

	var request PatchUserRequest
	if err := requests.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err, "failed to decode request:",
		)
		return
	}

	userPatch := convertUserPatchFromRequest(request)
	userDomain, err := c.UserService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err, "failed to patch user",
		)
		return
	}

	log.Debug("userDomain before response",
		zap.Int("id", userDomain.Id),
		zap.String("fullName", userDomain.FullName),
		zap.Any("phoneNumber", userDomain.PhoneNumber), // покажет nil или значение
	)

	response := PatchUserResponse(convertUserDTOFromDomain(userDomain))
	responseHandler.JsonResponse(response, http.StatusOK)
}

func convertUserPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.Phone_number.ToDomain(),
	)
}
