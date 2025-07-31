package mailutil

import (
	"fmt"
	"testing"

	"github.com/karosown/katool-go/container/stream"
)

func TestVerifyEduEmail(t *testing.T) {
	validator := NewEduEmailVerify()

	testEmails := []string{
		"student@harvard.edu",      // 美国
		"user@oxford.ac.uk",        // 英国
		"researcher@sydney.edu.au", // 澳大利亚
		"student@tsinghua.edu.cn",  // 中国
		"prof@u-tokyo.ac.jp",       // 日本
		"student@nus.edu.sg",       // 新加坡
		"user@iitd.ac.in",          // 印度
		"student@uct.ac.za",        // 南非
		"fake@university.com",      // 非教育邮箱
		"fake@service.liberty.edu", // 非教育邮箱
		"fuck@.vatican.va",
	}

	stream.Cast[CompleteEmailInfo](stream.ToStream(&testEmails).Map(func(i string) any {
		return validator.GetCompleteInfo(i)
	})).ForEach(func(info CompleteEmailInfo) {
		fmt.Printf("邮箱: %s\n", info.Email)
		fmt.Printf("  教育邮箱: %t\n", info.IsEducation)
		fmt.Printf("  国家: %s\n", info.Country)
		fmt.Printf("  机构: %s\n", info.Institution)
		fmt.Println()
	})

}
