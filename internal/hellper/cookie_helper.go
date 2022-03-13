package hellper

//func GetToken(cookies []*http.Cookie) string {
//	for _, cookie := range cookies {
//		// пробуем получить значение uid и hash из куки
//		if cookie.Name == CookieName {
//			parts := strings.Split(cookie.Value, ":")
//			if len(parts) != 2 {
//				// если в куки нет обоих параметров, то генерируем новый uid
//				return GenerateRandomString(uidLen)
//			}
//			uid, hash := parts[0], parts[1]
//			if checkHash(uid, hash) {
//				return uid
//			}
//		}
//	}
//	return GenerateRandomString(uidLen)
//}
