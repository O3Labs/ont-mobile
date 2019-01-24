package ontmobile

import (
  "bytes"
  "fmt"
  "math"
  "log"
  "encoding/hex"

  "github.com/ontio/ontology/cmd/utils"
  "github.com/ontio/ontology/common"
  httpcom "github.com/ontio/ontology/http/base/common"
)

func buildParameters(args []Parameter) ([]interface{}){
  var parameters []interface{}
  var err error

  for _, element := range args {
    var t = element.Type
    var v = element.Value
    var p interface{}
    if t == Address {
      p, err = common.AddressFromBase58(v.(string))
      if err != nil {
        log.Printf("Failed to convert string to address %s", err)
      }
    } else if t == String {
      p = v.(string)
    } else if t == Integer {
      p = v.(uint)
    } else if t == Fixed8 {
      p = uint64(RoundFixed(v.(float64), ONGDECIMALS) * float64(math.Pow10(ONGDECIMALS)))
    } else if t == Array {
      p = buildParameters(v.([]Parameter))
    }
    parameters = append(parameters, p)
  }

  return parameters
}

// BuildInvocationTransaction : creates a raw transaction
func BuildInvocationTransaction(contractHex string, operation string, args []Parameter, gasPrice uint, gasLimit uint, wif string) (string, error) {
  var contractAddress common.Address
  contractAddress, err := common.AddressFromHexString(contractHex)
  if err != nil {
    return "", fmt.Errorf("[Invalid contract hash error: %s]", err)
  }

  signer := ONTAccountWithWIF(wif).account
  parameters := buildParameters(args)
  params := []interface{}{operation, parameters}

  tx, err := httpcom.NewNeovmInvokeTransaction(uint64(gasPrice), uint64(gasLimit), contractAddress, params)
  if err != nil {
      log.Printf("NewNeovmInvokeTransaction error:%s", err)
      return "", err
  }

  tx.Payer = signer.Address

  err = utils.SignTransaction(signer, tx)
  if err != nil {
    log.Printf("SignTransaction error: %s", err)
    return "", err
  }

  immutTx, err := tx.IntoImmutable()

  if err != nil {
    log.Printf("IntoImmutable error: %s", err)
    return "", err
  }

  var buffer bytes.Buffer
  err = immutTx.Serialize(&buffer)
  if err != nil {
    log.Printf("serialize error:%s", err)
    return "", err
  }

  txData := hex.EncodeToString(buffer.Bytes())
  return txData, nil
}
