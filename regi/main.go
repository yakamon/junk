package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	ProductTypeRegular = iota
	ProductTypeQuantitative
)

func getProductMap(stdin *bufio.Scanner) map[string]Product {
	productMap := map[string]Product{}

	stdin.Scan()
	N, _ := strconv.Atoi(stdin.Text())
	for i := 0; i < N; i++ {
		stdin.Scan()
		line := strings.Split(stdin.Text(), " ")

		switch id := line[0]; len(id) {
		case 5: // 量り売り商品
			weightPerHundredYen, _ := strconv.ParseFloat(line[1], 64)
			packageWeight, _ := strconv.ParseFloat(line[2], 64)
			allowableErrorWeight, _ := strconv.ParseFloat(line[3], 64)

			productMap[id] = &QuantitativeProductInfo{id, weightPerHundredYen, packageWeight, allowableErrorWeight}
		case 12: // 通常商品
			price, _ := strconv.Atoi(line[1])
			standardWeight, _ := strconv.ParseFloat(line[2], 64)
			allowableErrorWeight, _ := strconv.ParseFloat(line[3], 64)

			productMap[id] = &RegularProductInfo{id, price, standardWeight, allowableErrorWeight}
		}
	}

	return productMap
}

func getAccounts(stdin *bufio.Scanner) []*Account {
	var accounts []*Account

	var curAccount *Account
	var prevScan *ProductScan
	var prevTotalWeight float64
	for stdin.Scan() {
		switch line := strings.Split(stdin.Text(), " "); line[0] {
		case "start":
			curAccount = &Account{[]*ProductScan{}}
			prevScan = &ProductScan{"", 0}
			prevTotalWeight = 0
		case "end":
			curTotalWeight, _ := strconv.ParseFloat(line[1], 64)
			prevScan.Weight = curTotalWeight - prevTotalWeight
			accounts = append(accounts, curAccount)
		default:
			barCode := line[0]
			curTotalWeight, _ := strconv.ParseFloat(line[1], 64)

			prevScan.Weight = curTotalWeight - prevTotalWeight

			curScan := &ProductScan{barCode, 0}
			curAccount.ProductScans = append(curAccount.ProductScans, curScan)

			prevScan = curScan
			prevTotalWeight = curTotalWeight
		}
	}

	return accounts
}

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	productMap := getProductMap(stdin)
	for i, account := range getAccounts(stdin) {
		fmt.Println("account:", i)
		handleAccount(account, productMap)
	}
}

func handleAccount(account *Account, productMap map[string]Product) {
	priceSum := 0
	isInvalidCode := false
	isInvalidWeight := false

	// 会計を処理
	for _, scan := range account.ProductScans {
		// チェックサム検証
		checkSum, _ := strconv.Atoi(string(scan.BarCode[12]))
		var digitSum int
		for _, v := range scan.BarCode[:12] {
			n, _ := strconv.Atoi(string(v))
			digitSum += n
		}
		if digitSum%10 != checkSum {
			isInvalidCode = true
		}

		// 商品の存在確認と取得
		var price int
		switch scan.getProductType() {
		case ProductTypeRegular:
			productID := scan.BarCode[:12]
			product, exists := productMap[productID]
			if !exists {
				isInvalidCode = true
			}
			// 重量検証
			if regularProduct := product.(*RegularProductInfo); !regularProduct.isValidWeight(scan.Weight) {
				isInvalidWeight = true
			} else {
				price = regularProduct.getPrice()
			}
		case ProductTypeQuantitative:
			productID := scan.BarCode[2:7]
			product, exists := productMap[productID]
			if !exists {
				isInvalidCode = true
			}
			// 重量検証
			scanPrice, _ := strconv.Atoi(string(scan.BarCode[7:12]))
			if quantitativeProduct := product.(*QuantitativeProductInfo); !quantitativeProduct.isValidWeight(scan.Weight, scanPrice) {
				isInvalidWeight = true
			} else {
				price = scanPrice
			}
		}

		priceSum += price
	}

	// 結果を出力
	if isInvalidCode || isInvalidWeight {
		callStaff(isInvalidCode, isInvalidWeight)
		return
	}
	showPriceSum(priceSum)
}

type Account struct {
	ProductScans []*ProductScan
}

type ProductScan struct {
	BarCode string
	Weight  float64
}

func (ps *ProductScan) getProductType() int {
	if ps.BarCode[0:2] == "02" {
		return ProductTypeQuantitative
	}
	return ProductTypeRegular
}

type Product interface {
	getType() int
}

type RegularProductInfo struct {
	ID                   string
	Price                int
	StandardWeight       float64
	AllowableErrorWeight float64
}

func (rp *RegularProductInfo) getType() int {
	return ProductTypeRegular
}

func (rp *RegularProductInfo) isValidWeight(weight float64) bool {
	return math.Abs(rp.StandardWeight-weight) <= rp.AllowableErrorWeight
}

func (rp *RegularProductInfo) getPrice() int {
	return rp.Price
}

type QuantitativeProductInfo struct {
	ID                   string
	WeightPerHundredYen  float64
	PackageWeight        float64
	AllowableErrorWeight float64
}

func (qp *QuantitativeProductInfo) getType() int {
	return ProductTypeQuantitative
}

func (qp *QuantitativeProductInfo) isValidWeight(weight float64, price int) bool {
	standardWeight := qp.WeightPerHundredYen*float64(price)/100 + qp.PackageWeight
	return math.Abs(standardWeight-weight) <= qp.AllowableErrorWeight
}

func callStaff(isInvalidCode, isInvalidWeight bool) {
	message := "staff call:"
	if isInvalidCode {
		message += " 1"
	}
	if isInvalidWeight {
		message += " 2"
	}
	fmt.Println(message)
}

func showPriceSum(priceSum int) {
	fmt.Println(priceSum)
}
