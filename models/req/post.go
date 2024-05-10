package req

import "mime/multipart"

type PostCreate struct {
	Caption       string                  `json:"caption" validate:"required,min=1,max=100"`
	Upload        []*multipart.FileHeader `json:"upload" validate:"required"`
	IsHideLike    bool                    `json:"is_hide_like"`
	IsHideComment bool                    `json:"is_hide_comment"`
}
