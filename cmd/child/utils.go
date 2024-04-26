package child

import (
	"strconv"
	"strings"

	"github.com/nerijusro/scootinAboot/types"
)

func getScootersByAreaUrl(basePath string, params types.GetScootersQueryParameters) string {
	var url strings.Builder
	url.WriteString(basePath)
	url.WriteString("/client/scooters?availability=")
	url.WriteString(string(params.Availability))
	url.WriteString("&x1=")
	url.WriteString(strconv.FormatFloat(params.X1, 'f', -1, 64))
	url.WriteString("&x2=")
	url.WriteString(strconv.FormatFloat(params.X2, 'f', -1, 64))
	url.WriteString("&y1=")
	url.WriteString(strconv.FormatFloat(params.Y1, 'f', -1, 64))
	url.WriteString("&y2=")
	url.WriteString(strconv.FormatFloat(params.Y2, 'f', -1, 64))

	return url.String()
}
