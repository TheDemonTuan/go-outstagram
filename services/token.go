package services

import (
	"errors"
	"outstagram/common"
	"outstagram/models/entity"
	"time"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (s *TokenService) SaveRefreshToken(userID, refreshToken string, expiredAt time.Time) error {
	newRefreshToken := entity.Token{
		UserID:       userID,
		RefreshToken: refreshToken,
		ExpiredAt:    expiredAt,
	}

	if err := common.DBConn.Create(&newRefreshToken).Error; err != nil {
		return errors.New("error while saving refresh token")
	}

	return nil
}

func (s *TokenService) GetRefreshToken(userID string) (entity.Token, error) {
	var token entity.Token

	if err := common.DBConn.Where("user_id = ? AND active = ?", userID, true).First(&token).Error; err != nil {
		return entity.Token{}, errors.New("refresh token not found")
	}

	if token.ExpiredAt.Before(time.Now()) {
		return entity.Token{}, errors.New("refresh token expired")
	}

	return token, nil
}

func (s *TokenService) GetRefreshTokenByToken(refreshToken string) (entity.Token, error) {
	var token entity.Token

	if err := common.DBConn.Where("refresh_token = ? AND active = ?", refreshToken, true).First(&token).Error; err != nil {
		return entity.Token{}, errors.New("refresh token not found")
	}

	if token.ExpiredAt.Before(time.Now()) {
		return entity.Token{}, errors.New("refresh token expired")
	}

	return token, nil
}

func (s *TokenService) DeleteRefreshTokenByUserID(userID string) error {
	if err := common.DBConn.Where("user_id = ?", userID).Delete(&entity.Token{}).Error; err != nil {
		return errors.New("error while deleting refresh token")
	}

	return nil
}

func (s *TokenService) DeleteRefreshTokenByToken(userID, refreshToken string) error {
	if err := common.DBConn.Where("refresh_token = ? AND user_id = ?", refreshToken, userID).Delete(&entity.Token{}).Error; err != nil {
		return errors.New("error while deleting refresh token")
	}

	return nil
}

func (s *TokenService) RevokeRefreshToken(userID string) error {
	if err := common.DBConn.Model(&entity.Token{}).Where("user_id = ?", userID).Update("active", false).Error; err != nil {
		return errors.New("error while revoking refresh token")
	}

	return nil
}

func (s *TokenService) RevokeAllRefreshToken(userID string) error {
	if err := common.DBConn.Model(&entity.Token{}).Where("user_id = ?", userID).Update("active", false).Error; err != nil {
		return errors.New("error while revoking all refresh token")
	}

	return nil
}

func (s *TokenService) CleanUpExpiredToken() error {
	if err := common.DBConn.Where("expired_at < ?", time.Now()).Delete(&entity.Token{}).Error; err != nil {
		return errors.New("error while cleaning up expired token")
	}

	return nil
}

func (s *TokenService) CleanUpAllToken() error {
	if err := common.DBConn.Delete(&entity.Token{}).Error; err != nil {
		return errors.New("error while cleaning up all token")
	}

	return nil
}

func (s *TokenService) CleanUpAllTokenByUserID(userID string) error {
	if err := common.DBConn.Where("user_id = ?", userID).Delete(&entity.Token{}).Error; err != nil {
		return errors.New("error while cleaning up all token by user id")
	}

	return nil
}
