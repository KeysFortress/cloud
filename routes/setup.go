package routes

import (
	"crypto/ed25519"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"leanmeal/api/dtos"
	"leanmeal/api/repositories"
	"leanmeal/api/utils"
)

type SetupController struct {
	accountRepository    repositories.Accounts
	accessKeysRepository repositories.AccessKeysRepository
	setupPath            string
}

func (s *SetupController) setup(ctx *gin.Context) {

	request := dtos.SetupRequest{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	s.accountRepository.Storage.Open()
	defer s.accountRepository.Storage.Close()

	initialSetup, err := s.accountRepository.IsEmpty()

	if err != nil || !initialSetup {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	code := utils.GenerateRandomString(32)
	uuid := uuid.New()

	ctx.Set("AuthRequest", dtos.StoredAuthRequest{
		Uuid:        uuid.String(),
		Code:        code,
		Name:        request.Email,
		ApprovedKey: request.PublicKey,
	})

	hostname, err := os.Hostname()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	result := "keysfortres://url=" + hostname + "&&setup=" + hostname + s.setupPath + "&&secret=" + code + "&&id=" + uuid.String()
	ctx.JSON(http.StatusOK, result)
}

func (s *SetupController) finish(ctx *gin.Context) {
	request := dtos.FinishAuthResponse{}
	if err := ctx.BindJSON(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	s.accountRepository.Storage.Open()
	defer s.accountRepository.Storage.Close()
	initialRequest := ctx.MustGet("AuthRequest").(dtos.StoredAuthRequest)
	publicKey, err := base64.StdEncoding.DecodeString(initialRequest.ApprovedKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	dedcodedMessege, err := base64.StdEncoding.DecodeString(initialRequest.Code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	decodedSignature, err := base64.StdEncoding.DecodeString(request.Signature)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	isValid := ed25519.Verify(publicKey, dedcodedMessege, decodedSignature)

	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	id, err := s.accountRepository.CreateAccount(&dtos.CreateAccountRequest{
		Email:     initialRequest.Name,
		Name:      "--",
		PublicKey: initialRequest.ApprovedKey,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	_, err = s.accessKeysRepository.Add(&id, &initialRequest.ApprovedKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Bad Request"})
		return
	}

	ctx.Set("AuthRequest", "")
	ctx.JSON(http.StatusOK, isValid)
}

func (s *SetupController) Init(r *gin.RouterGroup) {

	controller := r.Group("setup")
	controller.POST("init", s.setup)
	controller.POST("finish", s.finish)
}
