package nocomment

import (
	"bytes"
)

const (
	CODE                   = iota // 正常代码
	COMMENTS_MULTILINE            // 多行注释
	COMMENTS_SINGLELINE           // 单行注释
	COMMENTS_HTML                 // HTML注释
	BACKSLASH                     // 折行注释
	CODE_CHAR                     // 字符
	CHAR_ESCAPE_SEQUENCE          // 字符中的转义字符
	CODE_STRING                   // 字符串
	STRING_ESCAPE_SEQUENCE        // 字符串中的转义字符
)

type Stripper struct {
	// 删除 /* */ 风格的注释
	RemoveBlockComment bool
	// 删除 // 风格的注释
	RemoveLineComment bool
	// 删除 # 风格的注释
	RemoveShellComment bool
	// 删除 <!-- --> 风格的注释
	RemoveHtmlComment bool
	// 删除 -- 风格的注释
	RemoveSQLComment bool
}

func (stripper *Stripper) Clean(input []byte) []byte {
	var out bytes.Buffer
	state := CODE

	for i := 0; i < len(input); i++ {
		b := input[i]

		switch state {
		case CODE:
			switch b {
			case '/':
				index := i + 1
				if index < len(input) {
					// //
					if stripper.RemoveLineComment && input[index] == '/' {
						state = COMMENTS_SINGLELINE
						i += 1
						continue
					}
					// /*
					if stripper.RemoveBlockComment && input[index] == '*' {
						state = COMMENTS_MULTILINE
						i += 1
						continue
					}
				}
			case '-':
				index := i + 1
				if index < len(input) {
					// --
					if stripper.RemoveSQLComment && input[index] == '-' {
						state = COMMENTS_SINGLELINE
						i += 1
						continue
					}
				}
			case '#':
				if stripper.RemoveShellComment {
					state = COMMENTS_SINGLELINE
					continue
				}
			case '<':
				// <!--
				if i < len(input)-3 {
					if stripper.RemoveHtmlComment {
						if input[i+1] == '!' && input[i+2] == '-' && input[i+3] == '-' {
							state = COMMENTS_HTML
							i += 3
							continue
						}
					}
				}
			case '\'':
				state = CODE_CHAR
			case '"':
				state = CODE_STRING
			}

			out.WriteByte(b)
		case COMMENTS_MULTILINE:
			// */
			if b == '*' {
				index := i + 1
				if index < len(input) {
					if input[index] == '/' {
						state = CODE
						i += 1
					}
				}
			}
		case COMMENTS_SINGLELINE:
			if b == '\\' {
				state = BACKSLASH
			} else if b == '\n' || b == '\r' {
				out.WriteByte(b)
				state = CODE
			}
		case COMMENTS_HTML:
			// -->
			if b == '-' {
				if i < len(input)-2 {
					if input[i+1] == '-' && input[i+2] == '>' {
						state = CODE
						i += 2
					}
				}
			}
		case BACKSLASH:
			if b != '\\' && b != '\n' && b != '\r' {
				state = COMMENTS_SINGLELINE
			}
		case CODE_CHAR:
			out.WriteByte(b)
			if b == '\\' {
				state = CHAR_ESCAPE_SEQUENCE
			} else if b == '\'' {
				state = CODE
			}
		case CHAR_ESCAPE_SEQUENCE:
			out.WriteByte(b)
			state = CODE_CHAR
		case CODE_STRING:
			out.WriteByte(b)
			if b == '\\' {
				state = STRING_ESCAPE_SEQUENCE
			} else if b == '"' {
				state = CODE
			}
		case STRING_ESCAPE_SEQUENCE:
			out.WriteByte(b)
			state = CODE_STRING
		}
	}

	return out.Bytes()
}
