package wallet

import (
    "fmt"
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha256"
    "golang.org/x/crypto/ripemd160"
    "github.com/btcsuite/btcutil/base58"
)

type Wallet struct {
    privateKey        *ecdsa.PrivateKey
    publicKey         *ecdsa.PublicKey
    blockchainAddress string
}

func NewWallet() *Wallet {
    w := new(Wallet)
    privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
    w.privateKey = privateKey
    w.publicKey = &w.privateKey.PublicKey

    // 1 - Take the corresponding public key generated with it (33 bytes, 1 byte 0x02 (y-coord is even), and 32 bytes corresponding to X coordinate)

    // 2 - Perform SHA-256 hashing on the public key
    h2 := sha256.New()
    h2.Write(w.publicKey.X.Bytes())
    h2.Write(w.publicKey.Y.Bytes())
    digest2 := h2.Sum(nil)

    // 3 - Perform RIPEMD-160 hashing on the result of SHA-256
    h3 := ripemd160.New()
    h3.Write(digest2)
    digest3 := h3.Sum(nil)

    // 4 - Add version byte in front of RIPEMD-160 hash (0x00 for Main Network)
    vd4 := make([]byte, 21)
    vd4[0] = 0x00
    copy(vd4[1:], digest3[:])

    // 5 - Perform SHA-256 hash on the extended RIPEMD-160 result
    h5 := sha256.New()
    h5.Write(vd4)
    digest5 := h5.Sum(nil)

    // 6 - Perform SHA-256 hash on the result of the previous SHA-256 hash
    h6 := sha256.New()
    h6.Write(digest5)
    digest6 := h6.Sum(nil)
    
    // 7 - Take the first 4 bytes of the second SHA-256 hash. This is the address checksum
    chsum := digest6[:4]

    // 8 - Add the 4 checksum bytes from stage 7 at the end of extended RIPEMD-160 hash from stage 4. This is the 25-byte binary Bitcoin Address.
    dc8 := make([]byte, 25)
    copy(dc8[:21], vd4[:])
    copy(dc8[21:], chsum[:])

    // 9 - Convert the result from a byte string into a base58 string using Base58Check encoding. This is the most commonly used Bitcoin Address format
    address := base58.Encode(dc8)
    w.blockchainAddress = address
    return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
    return w.privateKey // struct
}

func (w *Wallet) PrivateKeyStr() string {
    return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
    return w.publicKey // struct
}

func (w *Wallet) PublicKeyStr() string {
    return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
    return w.blockchainAddress
}
