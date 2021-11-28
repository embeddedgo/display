// Code generated by mkcordic.go; DO NOT EDIT.

package math2d

const fracTh = 1

var cordicThs = [30]int32{
	1073741824, // 45.000000° = atan(1.000000000)
	633866811,  // 26.565051° = atan(0.500000000)
	334917815,  // 14.036243° = atan(0.250000000)
	170009512,  //  7.125016° = atan(0.125000000)
	85334662,   //  3.576334° = atan(0.062500000)
	42708931,   //  1.789911° = atan(0.031250000)
	21359677,   //  0.895174° = atan(0.015625000)
	10680490,   //  0.447614° = atan(0.007812500)
	5340327,    //  0.223811° = atan(0.003906250)
	2670173,    //  0.111906° = atan(0.001953125)
	1335088,    //  0.055953° = atan(0.000976562)
	667544,     //  0.027976° = atan(0.000488281)
	333772,     //  0.013988° = atan(0.000244141)
	166886,     //  0.006994° = atan(0.000122070)
	83443,      //  0.003497° = atan(0.000061035)
	41722,      //  0.001749° = atan(0.000030518)
	20861,      //  0.000874° = atan(0.000015259)
	10430,      //  0.000437° = atan(0.000007629)
	5215,       //  0.000219° = atan(0.000003815)
	2608,       //  0.000109° = atan(0.000001907)
	1304,       //  0.000055° = atan(0.000000954)
	652,        //  0.000027° = atan(0.000000477)
	326,        //  0.000014° = atan(0.000000238)
	163,        //  0.000007° = atan(0.000000119)
	81,         //  0.000003° = atan(0.000000060)
	41,         //  0.000002° = atan(0.000000030)
	20,         //  0.000001° = atan(0.000000015)
	10,         //  0.000000° = atan(0.000000007)
	5,          //  0.000000° = atan(0.000000004)
	3,          //  0.000000° = atan(0.000000002)
}

const (
	fracK   = 16
	cordicK = 39797 // 0.607253 * (1<<fracK)
)
