// Code generated by "stringer -type=tokenType"; DO NOT EDIT

package zonefile

import "fmt"

const _tokenType_name = "tokenErrortokenEOFtokenWhiteSpacetokenLeftParentokenRightParentokenCommenttokenItemtokenQuotedItemtokenNewline"

var _tokenType_index = [...]uint8{0, 10, 18, 33, 47, 62, 74, 83, 98, 110}

func (i tokenType) String() string {
	if i < 0 || i >= tokenType(len(_tokenType_index)-1) {
		return fmt.Sprintf("tokenType(%d)", i)
	}
	return _tokenType_name[_tokenType_index[i]:_tokenType_index[i+1]]
}