package nocomment

import (
	"bytes"
)

const (
	CODE                                = iota // 正常代码
	COMMENTS_MULTILINE                         // 多行注释
	COMMENTS_SINGLELINE                        // 单行注释
	COMMENTS_HTML                              // HTML注释
	BACKSLASH                                  // 折行注释
	SINGLEQUOTE_STRING                         // 单引号字符串
	SINGLEQUOTE_STRING_ESCAPE_SEQUENCE         // 单引号字符串中的转义字符
	DOUBLEQUOTES_STRING                        // 双引号字符串
	DOUBLEQUOTES_STRING_ESCAPE_SEQUENCE        // 双引号字符串中的转义字符
	BACKTICK_STRING                            // 反引号字符串
	BACKTICK_STRING_ESCAPE_SEQUENCE            // 反引号字符串中的转义字符
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
				if i+2 < len(input) {
					// --
					if stripper.RemoveSQLComment && input[i+1] == '-' {
						if input[i+2] == ' ' || input[i+2] == '\t' || input[i+2] == '\n' {
							state = COMMENTS_SINGLELINE
							i += 2
							continue
						}
						if input[i+2] == '\r' {
							state = COMMENTS_SINGLELINE
							i += 2
							if i+1 < len(input) && input[i+1] == '\n' {
								i += 1
							}
							continue
						}
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
				state = SINGLEQUOTE_STRING
			case '"':
				state = DOUBLEQUOTES_STRING
			case '`':
				state = BACKTICK_STRING
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
		case SINGLEQUOTE_STRING:
			out.WriteByte(b)
			if b == '\\' {
				state = SINGLEQUOTE_STRING_ESCAPE_SEQUENCE
			} else if b == '\'' {
				state = CODE
			}
		case SINGLEQUOTE_STRING_ESCAPE_SEQUENCE:
			out.WriteByte(b)
			state = SINGLEQUOTE_STRING
		case DOUBLEQUOTES_STRING:
			out.WriteByte(b)
			if b == '\\' {
				state = DOUBLEQUOTES_STRING_ESCAPE_SEQUENCE
			} else if b == '"' {
				state = CODE
			}
		case DOUBLEQUOTES_STRING_ESCAPE_SEQUENCE:
			out.WriteByte(b)
			state = DOUBLEQUOTES_STRING
		case BACKTICK_STRING:
			out.WriteByte(b)
			if b == '`' {
				state = BACKTICK_STRING_ESCAPE_SEQUENCE
			}
		case BACKTICK_STRING_ESCAPE_SEQUENCE:
			out.WriteByte(b)
			if b == '`' {
				state = BACKTICK_STRING
			} else {
				state = CODE
			}
		}
	}

	return out.Bytes()
}
