package options

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/imgproxy/imgproxy/v3/config"
	"github.com/imgproxy/imgproxy/v3/imagetype"
)

type ProcessingOptionsTestSuite struct{ suite.Suite }

func (s *ProcessingOptionsTestSuite) SetupTest() {
	config.Reset()
	// Reset presets
	presets = make(map[string]urlOptions)
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URL() {
	originURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/size:100:100/%s.png", base64.RawURLEncoding.EncodeToString([]byte(originURL)))
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URLWithoutExtension() {
	originURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/size:100:100/%s", base64.RawURLEncoding.EncodeToString([]byte(originURL)))
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.Unknown, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URLWithBase() {
	config.BaseURL = "http://images.dev/"

	originURL := "lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/size:100:100/%s.png", base64.RawURLEncoding.EncodeToString([]byte(originURL)))
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), fmt.Sprintf("%s%s", config.BaseURL, originURL), imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URLWithReplacement() {
	config.URLReplacements = []config.URLReplacement{
		{Regexp: regexp.MustCompile("^test://([^/]*)/"), Replacement: "test2://images.dev/${1}/dolor/"},
		{Regexp: regexp.MustCompile("^test2://"), Replacement: "http://"},
	}

	originURL := "test://lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/size:100:100/%s.png", base64.RawURLEncoding.EncodeToString([]byte(originURL)))
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), "http://images.dev/lorem/dolor/ipsum.jpg?param=value", imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURL() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	path := fmt.Sprintf("/size:100:100/plain/%s@png", originURL)
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURLWithoutExtension() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	path := fmt.Sprintf("/size:100:100/plain/%s", originURL)

	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.Unknown, po.Format)
}
func (s *ProcessingOptionsTestSuite) TestParsePlainURLEscaped() {
	originURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/size:100:100/plain/%s@png", url.PathEscape(originURL))
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURLWithBase() {
	config.BaseURL = "http://images.dev/"

	originURL := "lorem/ipsum.jpg"
	path := fmt.Sprintf("/size:100:100/plain/%s@png", originURL)
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), fmt.Sprintf("%s%s", config.BaseURL, originURL), imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURLWithReplacement() {
	config.URLReplacements = []config.URLReplacement{
		{Regexp: regexp.MustCompile("^test://([^/]*)/"), Replacement: "test2://images.dev/${1}/dolor/"},
		{Regexp: regexp.MustCompile("^test2://"), Replacement: "http://"},
	}

	originURL := "test://lorem/ipsum.jpg"
	path := fmt.Sprintf("/size:100:100/plain/%s@png", originURL)
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), "http://images.dev/lorem/dolor/ipsum.jpg", imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParsePlainURLEscapedWithBase() {
	config.BaseURL = "http://images.dev/"

	originURL := "lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/size:100:100/plain/%s@png", url.PathEscape(originURL))
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), fmt.Sprintf("%s%s", config.BaseURL, originURL), imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

// func (s *ProcessingOptionsTestSuite) TestParseURLAllowedSource() {
// 	config.AllowedSources = []string{"local://", "http://images.dev/"}

// 	path := "/plain/http://images.dev/lorem/ipsum.jpg"
// 	_, _, err := ParsePath(path, make(http.Header))

// 	require.Nil(s.T(), err)
// }

// func (s *ProcessingOptionsTestSuite) TestParseURLNotAllowedSource() {
// 	config.AllowedSources = []string{"local://", "http://images.dev/"}

// 	path := "/plain/s3://images/lorem/ipsum.jpg"
// 	_, _, err := ParsePath(path, make(http.Header))

// 	require.Error(s.T(), err)
// }

func (s *ProcessingOptionsTestSuite) TestParsePathFormat() {
	path := "/format:webp/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), imagetype.WEBP, po.Format)
}

func (s *ProcessingOptionsTestSuite) TestParsePathResize() {
	path := "/resize:fill:100:200:1/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), ResizeFill, po.ResizingType)
	require.Equal(s.T(), 100, po.Width)
	require.Equal(s.T(), 200, po.Height)
	require.True(s.T(), po.Enlarge)
}

func (s *ProcessingOptionsTestSuite) TestParsePathResizingType() {
	path := "/resizing_type:fill/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), ResizeFill, po.ResizingType)
}

func (s *ProcessingOptionsTestSuite) TestParsePathSize() {
	path := "/size:100:200:1/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 100, po.Width)
	require.Equal(s.T(), 200, po.Height)
	require.True(s.T(), po.Enlarge)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWidth() {
	path := "/width:100/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 100, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathHeight() {
	path := "/height:100/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 100, po.Height)
}

func (s *ProcessingOptionsTestSuite) TestParsePathEnlarge() {
	path := "/enlarge:1/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Enlarge)
}

func (s *ProcessingOptionsTestSuite) TestParsePathExtend() {
	path := "/extend:1:so:10:20/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), true, po.Extend.Enabled)
	require.Equal(s.T(), GravitySouth, po.Extend.Gravity.Type)
	require.Equal(s.T(), 10.0, po.Extend.Gravity.X)
	require.Equal(s.T(), 20.0, po.Extend.Gravity.Y)
}

func (s *ProcessingOptionsTestSuite) TestParsePathGravity() {
	path := "/gravity:soea/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), GravitySouthEast, po.Gravity.Type)
}

func (s *ProcessingOptionsTestSuite) TestParsePathGravityFocuspoint() {
	path := "/gravity:fp:0.5:0.75/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), GravityFocusPoint, po.Gravity.Type)
	require.Equal(s.T(), 0.5, po.Gravity.X)
	require.Equal(s.T(), 0.75, po.Gravity.Y)
}

func (s *ProcessingOptionsTestSuite) TestParsePathQuality() {
	path := "/quality:55/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 55, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParsePathBackground() {
	path := "/background:128:129:130/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Flatten)
	require.Equal(s.T(), uint8(128), po.Background.Color.R)
	require.Equal(s.T(), uint8(129), po.Background.Color.G)
	require.Equal(s.T(), uint8(130), po.Background.Color.B)
}

func (s *ProcessingOptionsTestSuite) TestParsePathBackgroundHex() {
	path := "/background:ffddee/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Flatten)
	require.Equal(s.T(), uint8(0xff), po.Background.Color.R)
	require.Equal(s.T(), uint8(0xdd), po.Background.Color.G)
	require.Equal(s.T(), uint8(0xee), po.Background.Color.B)
}

func (s *ProcessingOptionsTestSuite) TestParsePathAdvancedBackgroundEffect() {
	path := "/background:blur/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Flatten)
	require.Equal(s.T(), "blur", po.Background.Effect)
}

func (s *ProcessingOptionsTestSuite) TestParsePathBackgroundDisable() {
	path := "/background:fff/background:/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.False(s.T(), po.Flatten)
}

func (s *ProcessingOptionsTestSuite) TestParsePathBlur() {
	path := "/blur:0.2/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), float32(0.2), po.Blur)
}

func (s *ProcessingOptionsTestSuite) TestParsePathSharpen() {
	path := "/sharpen:0.2/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), float32(0.2), po.Sharpen)
}
func (s *ProcessingOptionsTestSuite) TestParsePathDpr() {
	path := "/dpr:2/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 2.0, po.Dpr)
}
func (s *ProcessingOptionsTestSuite) TestParsePathWatermark() {
	path := "/watermark:0.5:soea:10:20:0.6/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Watermark.Enabled)
	require.Equal(s.T(), GravitySouthEast, po.Watermark.Gravity.Type)
	require.Equal(s.T(), 10.0, po.Watermark.Gravity.X)
	require.Equal(s.T(), 20.0, po.Watermark.Gravity.Y)
	require.Equal(s.T(), 0.6, po.Watermark.Scale)
}

func (s *ProcessingOptionsTestSuite) TestParsePathPreset() {
	presets["test1"] = urlOptions{
		urlOption{Name: "resizing_type", Args: []string{"fill"}},
	}

	presets["test2"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	path := "/preset:test1:test2/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), ResizeFill, po.ResizingType)
	require.Equal(s.T(), float32(0.2), po.Blur)
	require.Equal(s.T(), 50, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParsePathPresetDefault() {
	presets["default"] = urlOptions{
		urlOption{Name: "resizing_type", Args: []string{"fill"}},
		urlOption{Name: "blur", Args: []string{"0.2"}},
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	path := "/quality:70/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), ResizeFill, po.ResizingType)
	require.Equal(s.T(), float32(0.2), po.Blur)
	require.Equal(s.T(), 70, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParsePathPresetLoopDetection() {
	presets["test1"] = urlOptions{
		urlOption{Name: "resizing_type", Args: []string{"fill"}},
	}

	presets["test2"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	path := "/preset:test1:test2:test1/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.ElementsMatch(s.T(), po.UsedPresets, []string{"test1", "test2"})
}

func (s *ProcessingOptionsTestSuite) TestParsePathCachebuster() {
	path := "/cachebuster:123/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), "123", po.CacheBuster)
}

func (s *ProcessingOptionsTestSuite) TestParsePathStripMetadata() {
	path := "/strip_metadata:true/plain/http://images.dev/lorem/ipsum.jpg"
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.StripMetadata)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWebpDetection() {
	config.EnableWebpDetection = true

	path := "/plain/http://images.dev/lorem/ipsum.jpg"
	headers := http.Header{"Accept": []string{"image/webp"}}
	po, _, err := ParsePath(path, headers)

	require.Nil(s.T(), err)

	require.Equal(s.T(), true, po.PreferWebP)
	require.Equal(s.T(), false, po.EnforceWebP)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWebpEnforce() {
	config.EnforceWebp = true

	path := "/plain/http://images.dev/lorem/ipsum.jpg@png"
	headers := http.Header{"Accept": []string{"image/webp"}}
	po, _, err := ParsePath(path, headers)

	require.Nil(s.T(), err)

	require.Equal(s.T(), true, po.PreferWebP)
	require.Equal(s.T(), true, po.EnforceWebP)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWidthHeader() {
	config.EnableClientHints = true

	path := "/plain/http://images.dev/lorem/ipsum.jpg@png"
	headers := http.Header{"Width": []string{"100"}}
	po, _, err := ParsePath(path, headers)

	require.Nil(s.T(), err)

	require.Equal(s.T(), 100, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWidthHeaderDisabled() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg@png"
	headers := http.Header{"Width": []string{"100"}}
	po, _, err := ParsePath(path, headers)

	require.Nil(s.T(), err)

	require.Equal(s.T(), 0, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathWidthHeaderRedefine() {
	config.EnableClientHints = true

	path := "/width:150/plain/http://images.dev/lorem/ipsum.jpg@png"
	headers := http.Header{"Width": []string{"100"}}
	po, _, err := ParsePath(path, headers)

	require.Nil(s.T(), err)

	require.Equal(s.T(), 150, po.Width)
}

func (s *ProcessingOptionsTestSuite) TestParsePathDprHeader() {
	config.EnableClientHints = true

	path := "/plain/http://images.dev/lorem/ipsum.jpg@png"
	headers := http.Header{"Dpr": []string{"2"}}
	po, _, err := ParsePath(path, headers)

	require.Nil(s.T(), err)

	require.Equal(s.T(), 2.0, po.Dpr)
}

func (s *ProcessingOptionsTestSuite) TestParsePathDprHeaderDisabled() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg@png"
	headers := http.Header{"Dpr": []string{"2"}}
	po, _, err := ParsePath(path, headers)

	require.Nil(s.T(), err)

	require.Equal(s.T(), 1.0, po.Dpr)
}

// func (s *ProcessingOptionsTestSuite) TestParsePathSigned() {
// 	config.Keys = [][]byte{[]byte("test-key")}
// 	config.Salts = [][]byte{[]byte("test-salt")}

// 	path := "/HcvNognEV1bW6f8zRqxNYuOkV0IUf1xloRb57CzbT4g/width:150/plain/http://images.dev/lorem/ipsum.jpg@png"
// 	_, _, err := ParsePath(path, make(http.Header))

// 	require.Nil(s.T(), err)
// }

// func (s *ProcessingOptionsTestSuite) TestParsePathSignedInvalid() {
// 	config.Keys = [][]byte{[]byte("test-key")}
// 	config.Salts = [][]byte{[]byte("test-salt")}

// 	path := "/unsafe/width:150/plain/http://images.dev/lorem/ipsum.jpg@png"
// 	_, _, err := ParsePath(path, make(http.Header))

// 	require.Error(s.T(), err)
// 	require.Equal(s.T(), signature.ErrInvalidSignature.Error(), err.Error())
// }

func (s *ProcessingOptionsTestSuite) TestParsePathOnlyPresets() {
	config.OnlyPresets = true
	presets["test1"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
	}
	presets["test2"] = urlOptions{
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	path := "/test1:test2/plain/http://images.dev/lorem/ipsum.jpg"

	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), float32(0.2), po.Blur)
	require.Equal(s.T(), 50, po.Quality)
}

func (s *ProcessingOptionsTestSuite) TestParseSkipProcessing() {
	path := "/skp:jpg:png/plain/http://images.dev/lorem/ipsum.jpg"

	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), []imagetype.Type{imagetype.JPEG, imagetype.PNG}, po.SkipProcessingFormats)
}

func (s *ProcessingOptionsTestSuite) TestParseSkipProcessingInvalid() {
	path := "/skp:jpg:png:bad_format/plain/http://images.dev/lorem/ipsum.jpg"

	_, _, err := ParsePath(path, make(http.Header))

	require.Error(s.T(), err)
	require.Equal(s.T(), "Invalid image format in skip processing: bad_format", err.Error())
}

func (s *ProcessingOptionsTestSuite) TestParseExpires() {
	path := "/exp:32503669200/plain/http://images.dev/lorem/ipsum.jpg"
	_, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
}

func (s *ProcessingOptionsTestSuite) TestParseExpiresExpired() {
	path := "/exp:1609448400/plain/http://images.dev/lorem/ipsum.jpg"
	_, _, err := ParsePath(path, make(http.Header))

	require.Error(s.T(), err)
	require.Equal(s.T(), errExpiredURL.Error(), err.Error())
}

func (s *ProcessingOptionsTestSuite) TestParseBase64URLOnlyPresets() {
	config.OnlyPresets = true
	presets["test1"] = urlOptions{
		urlOption{Name: "blur", Args: []string{"0.2"}},
	}
	presets["test2"] = urlOptions{
		urlOption{Name: "quality", Args: []string{"50"}},
	}

	originURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/test1:test2/%s.png", base64.RawURLEncoding.EncodeToString([]byte(originURL)))

	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), float32(0.2), po.Blur)
	require.Equal(s.T(), 50, po.Quality)
	require.Equal(s.T(), originURL, imageURL)
}

const nineColorPalette = "bk:12.79:5.21:-12.53:y:63.41:-12.92:68.53:bl:28.63:5.38:-39.49:" +
	"r:26.32:38.07:26.51:g:30.21:-22:5.54:w:66.16:-3.5:-1.98:" +
	"or:36.06:36.04:36.13:br:46.37:21.06:46.76:y2:68.07:-13.74:39.3"

const sixColorPalette = "bk:8.88:9.08:-12.38:y:60.02:-10.87:64.97:w:61.55:-3.11:-0.93:" +
	"bl:26.91:3.06:-34.11:r:23.76:38.93:28.27:g:31:-18.73:9.4"

func (s *ProcessingOptionsTestSuite) TestParseDitherOpts07() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	path := fmt.Sprintf("/dither:fs:opts07:0.65:%s/plain/%s", nineColorPalette, originURL)
	po, imageURL, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.True(s.T(), po.Dither.OptionsSet07)
	require.False(s.T(), po.Dither.OptionsSet06)
	require.Equal(s.T(), 0.65, po.Dither.SwapYellowProb)
	require.Equal(s.T(), nineColorPalette, po.Dither.MeasuredPalette)
}

func (s *ProcessingOptionsTestSuite) TestParseDitherOpts07WithTrailingOption() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	// applyDitherOption's loop never visits the last argument, so co is there to
	// give sp a position the loop actually reaches
	path := fmt.Sprintf("/dither:fs:opts07:0.65:%s:sp:co/plain/%s", nineColorPalette, originURL)
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.True(s.T(), po.Dither.OptionsSet07)
	require.Equal(s.T(), nineColorPalette, po.Dither.MeasuredPalette)
	require.True(s.T(), po.Dither.SoftProof, "the palette must not swallow the option after it")
}

func (s *ProcessingOptionsTestSuite) TestParseDitherOpts07RejectsSixColorPalette() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	path := fmt.Sprintf("/dither:fs:opts07:0.65:%s/plain/%s", sixColorPalette, originURL)
	_, _, err := ParsePath(path, make(http.Header))

	require.Error(s.T(), err, "opts07 carries 9 colors; a 6 color palette belongs in opts06")
}

func (s *ProcessingOptionsTestSuite) TestParseDitherOpts06StillTakesSixColors() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	path := fmt.Sprintf("/dither:fs:opts06:0.65:%s/plain/%s", sixColorPalette, originURL)
	po, _, err := ParsePath(path, make(http.Header))

	require.Nil(s.T(), err)
	require.True(s.T(), po.Dither.OptionsSet06)
	require.False(s.T(), po.Dither.OptionsSet07)
	require.Equal(s.T(), sixColorPalette, po.Dither.MeasuredPalette)
}

func TestProcessingOptions(t *testing.T) {
	suite.Run(t, new(ProcessingOptionsTestSuite))
}
