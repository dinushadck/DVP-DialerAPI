package main

import (
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/fzzy/radix/redis"
	"strconv"
	"strings"
	"time"
)

func loadJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	return (jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			claims := token.Claims.(jwt.MapClaims)
			secretKey := fmt.Sprintf("token:iss:%s:%s", claims["iss"], claims["jti"])
			secret := SecurityGet(secretKey)
			if secret == "" {
				return nil, fmt.Errorf("Invalied 'iss' or 'jti' in JWT")
			}
			return []byte(secret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	}))
}

func FindScope(vs []interface{}, scope, action string) (bool, bool) {
	for _, v := range vs {
		scopeInfo := v.(map[string]interface{})
		if scopeInfo["resource"] == scope {
			return true, FindAction(scopeInfo["actions"].([]interface{}), action)
		}
	}
	return false, false
}

func FindAction(vs []interface{}, action string) bool {
	for _, v := range vs {
		if v == action {
			return true
		}
	}
	return false
}

func decodeJwtDialer(dvp DVP, funcScope, action string) (company, tenant int, iss string, veeryMessage Result) {
	tokenVals := strings.Split(dvp.Context.Request().Header.Get("authorization"), " ")
	internalAccessToken := dvp.Context.Request().Header.Get("companyinfo")
	if len(tokenVals) > 1 {
		token, err := jwt.Parse(tokenVals[1], func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			claims := token.Claims.(jwt.MapClaims)
			secretKey := fmt.Sprintf("token:iss:%s:%s", claims["iss"], claims["jti"])
			//fmt.Println(secretKey)
			secret := SecurityGet(secretKey)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			scopes := claims["scope"].([]interface{})
			scope, actions := FindScope(scopes, funcScope, action)
			//fmt.Println(scope, ":: ", actions)
			if scope && actions {
				iiss := claims["iss"]
				iss := iiss.(string)
				fmt.Println(iiss)
				fmt.Println(iss)
				if internalAccessToken != "" {
					ids := strings.Split(internalAccessToken, ":")
					if len(ids) == 2 {
						tenant, _ = strconv.Atoi(ids[0])
						company, _ = strconv.Atoi(ids[1])
						return company, tenant, iss, ResponseGenerator(true, "Read Companyinfo Success", "", "")
					} else {
						return 0, 0, iss, ResponseGenerator(false, "Invalied Companyinfo", "", "")
					}
				} else {
					iTenant := claims["tenant"]
					iCompany := claims["company"]
					if iTenant != nil && iCompany != nil {
						tenant := int(iTenant.(float64))
						company := int(iCompany.(float64))
						return company, tenant, iss, ResponseGenerator(true, "Read Companyinfo Success", "", "")
					} else {
						return 0, 0, iss, ResponseGenerator(false, "Invalid company or tenant", "", "")
					}
				}
			} else {
				return 0, 0, "", ResponseGenerator(false, "Invalid scopes", "", "")
			}
		} else {
			fmt.Println(err)
			return 0, 0, "", ResponseGenerator(false, "Invalid token", "", "")
		}
	} else {
		return 0, 0, "", ResponseGenerator(false, "Invalid token", "", "")
	}
}

//func validateCompanyTenant(dvp DVP) (company, tenant int) {
//	internalAccessToken := dvp.Context.Request().Header.Get("companyinfo")
//	if internalAccessToken != "" {
//		ids := strings.Split(internalAccessToken, ":")
//		if len(ids) == 2 {
//			tenant, _ = strconv.Atoi(ids[0])
//			company, _ = strconv.Atoi(ids[1])
//			return company, tenant
//		} else {
//			return 0, 0
//		}
//	} else {
//		user := context.Get(dvp.Context.Request(), "user")
//		if user != nil {
//			iTenant := user.(*jwt.Token).Claims["tenant"]
//			iCompany := user.(*jwt.Token).Claims["company"]
//			if iTenant != nil && iCompany != nil {
//				tenant := int(iTenant.(float64))
//				company := int(iCompany.(float64))
//				return company, tenant
//			} else {
//				dvp.RB().Write(ResponseGenerator(false, "Invalid company or tenant", "", ""))
//				return
//			}
//		} else {
//			dvp.RB().Write(ResponseGenerator(false, "User data not found in JWT", "", ""))
//			return
//		}
//	}
//}

func ResponseGenerator(isSuccess bool, customMessage, result, exception string) Result {
	res := Result{}
	res.IsSuccess = isSuccess
	res.CustomMessage = customMessage
	res.Exception = exception
	res.Result = result
	return res
}

func SecurityGet(key string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", securityIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)

	strObj, _ := client.Cmd("get", key).Str()
	//fmt.Println(strObj)
	return strObj
}
