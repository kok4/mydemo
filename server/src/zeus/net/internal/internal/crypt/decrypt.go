package crypt

// decryptData 解密算法
func DecryptData(buf []byte) []byte {
	buflen := len(buf)
	key := decrypt_key
	keylen := len(key)

	for i := 0; i < buflen; i++ {
		b := buf[i] ^ key[i%keylen]
		n := byte(i%7 + 1)                 //移位长度(1-7)
		buf[i] = (b >> n) | (b << (8 - n)) // 向右循环移位
	}
	return buf
}
