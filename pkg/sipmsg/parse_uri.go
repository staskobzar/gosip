
//line parse_uri.rl:1
// -*-go-*-
//
// SIP URI parser

package sipmsg

import (
	"fmt"
)

func ParseURI(data string) (*URI, error) {
	
//line parse_uri.rl:13
	
//line parse_uri.go:18
var _uri_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3, 
	1, 4, 1, 5, 1, 6, 2, 1, 
	0, 
}

var _uri_key_offsets []int16 = []int16{
	0, 0, 2, 4, 6, 9, 25, 38, 
	44, 50, 64, 70, 76, 83, 91, 98, 
	106, 112, 119, 121, 134, 140, 146, 159, 
	172, 178, 184, 199, 215, 221, 227, 233, 
	239, 254, 260, 268, 274, 282, 288, 296, 
	304, 312, 320, 328, 336, 343, 351, 359, 
	367, 369, 376, 385, 387, 390, 392, 395, 
	397, 400, 403, 404, 407, 408, 411, 412, 
	421, 430, 438, 446, 454, 462, 464, 470, 
	479, 488, 497, 499, 502, 505, 506, 507, 
	526, 544, 563, 580, 598, 612, 630, 636, 
	642, 648, 654, 671, 677, 683, 701, 719, 
	725, 731, 749, 767, 773, 779, 798, 804, 
	810, 816, 822, 841, 847, 853, 871, 888, 
	907, 924, 943, 960, 979, 998, 1017, 1036, 
	1055, 1074, 1075, 1086, 1095, 1099, 1103, 1107, 
	1111, 1113, 1129, 1144, 1159, 1170, 1181, 1192, 
	1195, 1214, 1231, 1247, 1263, 1279, 1295, 1311, 
	1329, 1348, 1367, 1385, 1404, 1422, 1441, 1460, 
}

var _uri_trans_keys []byte = []byte{
	83, 115, 73, 105, 80, 112, 58, 83, 
	115, 33, 37, 59, 61, 63, 91, 95, 
	126, 36, 47, 48, 57, 65, 90, 97, 
	122, 33, 37, 58, 61, 64, 95, 126, 
	36, 59, 63, 90, 97, 122, 48, 57, 
	65, 70, 97, 102, 48, 57, 65, 70, 
	97, 102, 33, 37, 61, 64, 95, 126, 
	36, 46, 48, 57, 65, 90, 97, 122, 
	48, 57, 65, 70, 97, 102, 48, 57, 
	65, 70, 97, 102, 91, 48, 57, 65, 
	90, 97, 122, 45, 46, 48, 57, 65, 
	90, 97, 122, 45, 48, 57, 65, 90, 
	97, 122, 45, 46, 48, 57, 65, 90, 
	97, 122, 48, 57, 65, 90, 97, 122, 
	45, 48, 57, 65, 90, 97, 122, 48, 
	57, 33, 37, 93, 95, 126, 36, 43, 
	45, 58, 65, 91, 97, 122, 48, 57, 
	65, 70, 97, 102, 48, 57, 65, 70, 
	97, 102, 33, 37, 93, 95, 126, 36, 
	43, 45, 58, 65, 91, 97, 122, 33, 
	37, 93, 95, 126, 36, 43, 45, 58, 
	65, 91, 97, 122, 48, 57, 65, 70, 
	97, 102, 48, 57, 65, 70, 97, 102, 
	33, 36, 37, 63, 93, 95, 126, 39, 
	43, 45, 58, 65, 91, 97, 122, 33, 
	36, 37, 61, 63, 93, 95, 126, 39, 
	43, 45, 58, 65, 91, 97, 122, 48, 
	57, 65, 70, 97, 102, 48, 57, 65, 
	70, 97, 102, 48, 57, 65, 70, 97, 
	102, 48, 57, 65, 70, 97, 102, 33, 
	36, 37, 63, 93, 95, 126, 39, 43, 
	45, 58, 65, 91, 97, 122, 48, 57, 
	65, 90, 97, 122, 45, 46, 48, 57, 
	65, 90, 97, 122, 48, 57, 65, 90, 
	97, 122, 45, 46, 48, 57, 65, 90, 
	97, 122, 48, 57, 65, 90, 97, 122, 
	45, 46, 48, 57, 65, 90, 97, 122, 
	45, 46, 48, 57, 65, 90, 97, 122, 
	45, 46, 48, 57, 65, 90, 97, 122, 
	45, 46, 48, 57, 65, 90, 97, 122, 
	45, 46, 48, 57, 65, 90, 97, 122, 
	45, 46, 48, 57, 65, 90, 97, 122, 
	58, 48, 57, 65, 70, 97, 102, 58, 
	93, 48, 57, 65, 70, 97, 102, 58, 
	93, 48, 57, 65, 70, 97, 102, 58, 
	93, 48, 57, 65, 70, 97, 102, 58, 
	93, 58, 48, 57, 65, 70, 97, 102, 
	46, 58, 93, 48, 57, 65, 70, 97, 
	102, 48, 57, 46, 48, 57, 48, 57, 
	46, 48, 57, 48, 57, 93, 48, 57, 
	93, 48, 57, 93, 46, 48, 57, 46, 
	46, 48, 57, 46, 46, 58, 93, 48, 
	57, 65, 70, 97, 102, 46, 58, 93, 
	48, 57, 65, 70, 97, 102, 58, 93, 
	48, 57, 65, 70, 97, 102, 58, 93, 
	48, 57, 65, 70, 97, 102, 58, 93, 
	48, 57, 65, 70, 97, 102, 58, 93, 
	48, 57, 65, 70, 97, 102, 58, 93, 
	48, 57, 65, 70, 97, 102, 46, 58, 
	93, 48, 57, 65, 70, 97, 102, 46, 
	58, 93, 48, 57, 65, 70, 97, 102, 
	46, 58, 93, 48, 57, 65, 70, 97, 
	102, 48, 57, 46, 48, 57, 46, 48, 
	57, 46, 58, 33, 37, 45, 46, 58, 
	59, 61, 63, 64, 95, 126, 36, 47, 
	48, 57, 65, 90, 97, 122, 33, 37, 
	45, 58, 59, 61, 63, 64, 95, 126, 
	36, 47, 48, 57, 65, 90, 97, 122, 
	33, 37, 45, 46, 58, 59, 61, 63, 
	64, 95, 126, 36, 47, 48, 57, 65, 
	90, 97, 122, 33, 37, 58, 59, 61, 
	63, 64, 95, 126, 36, 47, 48, 57, 
	65, 90, 97, 122, 33, 37, 45, 58, 
	59, 61, 63, 64, 95, 126, 36, 47, 
	48, 57, 65, 90, 97, 122, 33, 37, 
	61, 64, 95, 126, 36, 46, 48, 57, 
	65, 90, 97, 122, 33, 37, 44, 58, 
	59, 61, 63, 64, 91, 93, 95, 126, 
	36, 57, 65, 90, 97, 122, 48, 57, 
	65, 70, 97, 102, 48, 57, 65, 70, 
	97, 102, 48, 57, 65, 70, 97, 102, 
	48, 57, 65, 70, 97, 102, 33, 37, 
	44, 47, 58, 61, 64, 91, 93, 95, 
	126, 36, 57, 65, 90, 97, 122, 48, 
	57, 65, 70, 97, 102, 48, 57, 65, 
	70, 97, 102, 33, 37, 44, 58, 59, 
	61, 63, 64, 91, 93, 95, 126, 36, 
	57, 65, 90, 97, 122, 33, 37, 44, 
	58, 59, 61, 63, 64, 91, 93, 95, 
	126, 36, 57, 65, 90, 97, 122, 48, 
	57, 65, 70, 97, 102, 48, 57, 65, 
	70, 97, 102, 33, 37, 38, 44, 58, 
	59, 61, 64, 91, 93, 95, 126, 36, 
	57, 63, 90, 97, 122, 33, 37, 38, 
	44, 58, 59, 61, 64, 91, 93, 95, 
	126, 36, 57, 63, 90, 97, 122, 48, 
	57, 65, 70, 97, 102, 48, 57, 65, 
	70, 97, 102, 33, 37, 38, 44, 47, 
	58, 61, 63, 64, 91, 93, 95, 126, 
	36, 57, 65, 90, 97, 122, 48, 57, 
	65, 70, 97, 102, 48, 57, 65, 70, 
	97, 102, 48, 57, 65, 70, 97, 102, 
	48, 57, 65, 70, 97, 102, 33, 37, 
	38, 44, 47, 58, 61, 63, 64, 91, 
	93, 95, 126, 36, 57, 65, 90, 97, 
	122, 48, 57, 65, 70, 97, 102, 48, 
	57, 65, 70, 97, 102, 33, 37, 38, 
	44, 58, 59, 61, 64, 91, 93, 95, 
	126, 36, 57, 63, 90, 97, 122, 33, 
	37, 58, 59, 61, 63, 64, 95, 126, 
	36, 47, 48, 57, 65, 90, 97, 122, 
	33, 37, 45, 46, 58, 59, 61, 63, 
	64, 95, 126, 36, 47, 48, 57, 65, 
	90, 97, 122, 33, 37, 58, 59, 61, 
	63, 64, 95, 126, 36, 47, 48, 57, 
	65, 90, 97, 122, 33, 37, 45, 46, 
	58, 59, 61, 63, 64, 95, 126, 36, 
	47, 48, 57, 65, 90, 97, 122, 33, 
	37, 58, 59, 61, 63, 64, 95, 126, 
	36, 47, 48, 57, 65, 90, 97, 122, 
	33, 37, 45, 46, 58, 59, 61, 63, 
	64, 95, 126, 36, 47, 48, 57, 65, 
	90, 97, 122, 33, 37, 45, 46, 58, 
	59, 61, 63, 64, 95, 126, 36, 47, 
	48, 57, 65, 90, 97, 122, 33, 37, 
	45, 46, 58, 59, 61, 63, 64, 95, 
	126, 36, 47, 48, 57, 65, 90, 97, 
	122, 33, 37, 45, 46, 58, 59, 61, 
	63, 64, 95, 126, 36, 47, 48, 57, 
	65, 90, 97, 122, 33, 37, 45, 46, 
	58, 59, 61, 63, 64, 95, 126, 36, 
	47, 48, 57, 65, 90, 97, 122, 33, 
	37, 45, 46, 58, 59, 61, 63, 64, 
	95, 126, 36, 47, 48, 57, 65, 90, 
	97, 122, 58, 45, 46, 58, 59, 63, 
	48, 57, 65, 90, 97, 122, 58, 59, 
	63, 48, 57, 65, 90, 97, 122, 59, 
	63, 48, 57, 59, 63, 48, 57, 59, 
	63, 48, 57, 59, 63, 48, 57, 59, 
	63, 33, 37, 59, 61, 63, 93, 95, 
	126, 36, 43, 45, 58, 65, 91, 97, 
	122, 33, 37, 59, 63, 93, 95, 126, 
	36, 43, 45, 58, 65, 91, 97, 122, 
	33, 37, 38, 63, 93, 95, 126, 36, 
	43, 45, 58, 65, 91, 97, 122, 45, 
	46, 58, 59, 63, 48, 57, 65, 90, 
	97, 122, 45, 46, 58, 59, 63, 48, 
	57, 65, 90, 97, 122, 45, 46, 58, 
	59, 63, 48, 57, 65, 90, 97, 122, 
	58, 59, 63, 33, 37, 45, 46, 58, 
	59, 61, 63, 64, 95, 126, 36, 47, 
	48, 57, 65, 90, 97, 122, 33, 37, 
	58, 59, 61, 63, 64, 95, 126, 36, 
	47, 48, 57, 65, 90, 97, 122, 33, 
	37, 59, 61, 63, 64, 95, 126, 36, 
	46, 48, 57, 65, 90, 97, 122, 33, 
	37, 59, 61, 63, 64, 95, 126, 36, 
	46, 48, 57, 65, 90, 97, 122, 33, 
	37, 59, 61, 63, 64, 95, 126, 36, 
	46, 48, 57, 65, 90, 97, 122, 33, 
	37, 59, 61, 63, 64, 95, 126, 36, 
	46, 48, 57, 65, 90, 97, 122, 33, 
	37, 59, 61, 63, 64, 95, 126, 36, 
	46, 48, 57, 65, 90, 97, 122, 33, 
	37, 44, 58, 59, 61, 63, 64, 91, 
	93, 95, 126, 36, 57, 65, 90, 97, 
	122, 33, 37, 44, 47, 58, 59, 61, 
	63, 64, 91, 93, 95, 126, 36, 57, 
	65, 90, 97, 122, 33, 37, 44, 47, 
	58, 59, 61, 63, 64, 91, 93, 95, 
	126, 36, 57, 65, 90, 97, 122, 33, 
	37, 44, 58, 59, 61, 63, 64, 91, 
	93, 95, 126, 36, 57, 65, 90, 97, 
	122, 33, 37, 38, 44, 47, 58, 61, 
	63, 64, 91, 93, 95, 126, 36, 57, 
	65, 90, 97, 122, 33, 37, 38, 44, 
	58, 59, 61, 64, 91, 93, 95, 126, 
	36, 57, 63, 90, 97, 122, 33, 37, 
	45, 46, 58, 59, 61, 63, 64, 95, 
	126, 36, 47, 48, 57, 65, 90, 97, 
	122, 33, 37, 45, 46, 58, 59, 61, 
	63, 64, 95, 126, 36, 47, 48, 57, 
	65, 90, 97, 122, 33, 37, 45, 46, 
	58, 59, 61, 63, 64, 95, 126, 36, 
	47, 48, 57, 65, 90, 97, 122, 
}

var _uri_single_lengths []byte = []byte{
	0, 2, 2, 2, 3, 8, 7, 0, 
	0, 6, 0, 0, 1, 2, 1, 2, 
	0, 1, 0, 5, 0, 0, 5, 5, 
	0, 0, 7, 8, 0, 0, 0, 0, 
	7, 0, 2, 0, 2, 0, 2, 2, 
	2, 2, 2, 2, 1, 2, 2, 2, 
	2, 1, 3, 0, 1, 0, 1, 0, 
	1, 1, 1, 1, 1, 1, 1, 3, 
	3, 2, 2, 2, 2, 2, 0, 3, 
	3, 3, 0, 1, 1, 1, 1, 11, 
	10, 11, 9, 10, 6, 12, 0, 0, 
	0, 0, 11, 0, 0, 12, 12, 0, 
	0, 12, 12, 0, 0, 13, 0, 0, 
	0, 0, 13, 0, 0, 12, 9, 11, 
	9, 11, 9, 11, 11, 11, 11, 11, 
	11, 1, 5, 3, 2, 2, 2, 2, 
	2, 8, 7, 7, 5, 5, 5, 3, 
	11, 9, 8, 8, 8, 8, 8, 12, 
	13, 13, 12, 13, 12, 11, 11, 11, 
}

var _uri_range_lengths []byte = []byte{
	0, 0, 0, 0, 0, 4, 3, 3, 
	3, 4, 3, 3, 3, 3, 3, 3, 
	3, 3, 1, 4, 3, 3, 4, 4, 
	3, 3, 4, 4, 3, 3, 3, 3, 
	4, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	0, 3, 3, 1, 1, 1, 1, 1, 
	1, 1, 0, 1, 0, 1, 0, 3, 
	3, 3, 3, 3, 3, 0, 3, 3, 
	3, 3, 1, 1, 1, 0, 0, 4, 
	4, 4, 4, 4, 4, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 3, 3, 
	3, 3, 3, 3, 3, 3, 4, 4, 
	4, 4, 4, 4, 4, 4, 4, 4, 
	4, 0, 3, 3, 1, 1, 1, 1, 
	0, 4, 4, 4, 3, 3, 3, 0, 
	4, 4, 4, 4, 4, 4, 4, 3, 
	3, 3, 3, 3, 3, 4, 4, 4, 
}

var _uri_index_offsets []int16 = []int16{
	0, 0, 3, 6, 9, 13, 26, 37, 
	41, 45, 56, 60, 64, 69, 75, 80, 
	86, 90, 95, 97, 107, 111, 115, 125, 
	135, 139, 143, 155, 168, 172, 176, 180, 
	184, 196, 200, 206, 210, 216, 220, 226, 
	232, 238, 244, 250, 256, 261, 267, 273, 
	279, 282, 287, 294, 296, 299, 301, 304, 
	306, 309, 312, 314, 317, 319, 322, 324, 
	331, 338, 344, 350, 356, 362, 365, 369, 
	376, 383, 390, 392, 395, 398, 400, 402, 
	418, 433, 449, 463, 478, 489, 505, 509, 
	513, 517, 521, 536, 540, 544, 560, 576, 
	580, 584, 600, 616, 620, 624, 641, 645, 
	649, 653, 657, 674, 678, 682, 698, 712, 
	728, 742, 758, 772, 788, 804, 820, 836, 
	852, 868, 870, 879, 886, 890, 894, 898, 
	902, 905, 918, 930, 942, 951, 960, 969, 
	973, 989, 1003, 1016, 1029, 1042, 1055, 1068, 
	1084, 1101, 1118, 1134, 1151, 1167, 1183, 1199, 
}

var _uri_indicies []byte = []byte{
	0, 0, 1, 2, 2, 1, 3, 3, 
	1, 4, 5, 5, 1, 6, 7, 6, 
	6, 6, 10, 6, 6, 6, 8, 9, 
	9, 1, 11, 12, 13, 11, 14, 11, 
	11, 11, 11, 11, 1, 15, 15, 15, 
	1, 11, 11, 11, 1, 13, 16, 13, 
	14, 13, 13, 13, 13, 13, 13, 1, 
	17, 17, 17, 1, 13, 13, 13, 1, 
	10, 18, 19, 19, 1, 20, 21, 22, 
	23, 23, 1, 20, 23, 23, 23, 1, 
	20, 24, 23, 23, 23, 1, 23, 25, 
	25, 1, 26, 25, 25, 25, 1, 27, 
	1, 28, 29, 28, 28, 28, 28, 28, 
	28, 28, 1, 30, 30, 30, 1, 31, 
	31, 31, 1, 31, 32, 31, 31, 31, 
	31, 31, 31, 31, 1, 33, 34, 33, 
	33, 33, 33, 33, 33, 33, 1, 35, 
	35, 35, 1, 33, 33, 33, 1, 36, 
	36, 37, 36, 36, 36, 36, 36, 36, 
	36, 36, 1, 38, 38, 39, 40, 38, 
	38, 38, 38, 38, 38, 38, 38, 1, 
	41, 41, 41, 1, 38, 38, 38, 1, 
	42, 42, 42, 1, 40, 40, 40, 1, 
	38, 38, 39, 38, 38, 38, 38, 38, 
	38, 38, 38, 1, 43, 25, 25, 1, 
	20, 44, 45, 23, 23, 1, 46, 25, 
	25, 1, 20, 47, 48, 23, 23, 1, 
	49, 25, 25, 1, 20, 47, 50, 23, 
	23, 1, 20, 47, 23, 23, 23, 1, 
	20, 44, 51, 23, 23, 1, 20, 44, 
	23, 23, 23, 1, 20, 21, 52, 23, 
	23, 1, 20, 21, 23, 23, 23, 1, 
	54, 53, 53, 53, 1, 56, 57, 55, 
	55, 55, 1, 56, 57, 58, 58, 58, 
	1, 56, 57, 59, 59, 59, 1, 56, 
	57, 1, 61, 60, 53, 53, 1, 62, 
	56, 57, 63, 55, 55, 1, 64, 1, 
	65, 66, 1, 67, 1, 68, 69, 1, 
	70, 1, 57, 71, 1, 57, 72, 1, 
	57, 1, 68, 73, 1, 68, 1, 65, 
	74, 1, 65, 1, 62, 56, 57, 75, 
	58, 58, 1, 62, 56, 57, 59, 59, 
	59, 1, 77, 57, 76, 76, 76, 1, 
	79, 57, 78, 78, 78, 1, 79, 57, 
	80, 80, 80, 1, 79, 57, 81, 81, 
	81, 1, 79, 57, 1, 82, 76, 76, 
	1, 62, 79, 57, 83, 78, 78, 1, 
	62, 79, 57, 84, 80, 80, 1, 62, 
	79, 57, 81, 81, 81, 1, 85, 1, 
	62, 86, 1, 62, 87, 1, 62, 1, 
	61, 1, 11, 12, 88, 89, 13, 11, 
	11, 11, 14, 11, 11, 11, 90, 91, 
	91, 1, 11, 12, 88, 13, 11, 11, 
	11, 14, 11, 11, 11, 91, 91, 91, 
	1, 11, 12, 88, 92, 13, 11, 11, 
	11, 14, 11, 11, 11, 91, 91, 91, 
	1, 11, 12, 13, 11, 11, 11, 14, 
	11, 11, 11, 91, 93, 93, 1, 11, 
	12, 94, 13, 11, 11, 11, 14, 11, 
	11, 11, 93, 93, 93, 1, 13, 16, 
	13, 14, 13, 13, 13, 95, 13, 13, 
	1, 96, 97, 11, 98, 11, 11, 11, 
	14, 28, 28, 96, 96, 96, 96, 96, 
	1, 99, 99, 99, 1, 100, 100, 100, 
	1, 101, 101, 101, 1, 102, 102, 102, 
	1, 103, 104, 13, 33, 33, 13, 14, 
	33, 33, 103, 103, 103, 103, 103, 1, 
	105, 105, 105, 1, 103, 103, 103, 1, 
	100, 106, 11, 102, 11, 11, 11, 14, 
	31, 31, 100, 100, 100, 100, 100, 1, 
	107, 108, 11, 103, 11, 11, 11, 14, 
	33, 33, 107, 107, 107, 107, 107, 1, 
	109, 109, 109, 1, 107, 107, 107, 1, 
	110, 111, 11, 11, 112, 11, 11, 14, 
	36, 36, 110, 110, 110, 110, 110, 1, 
	113, 114, 11, 11, 115, 11, 116, 14, 
	38, 38, 113, 113, 113, 113, 113, 1, 
	117, 117, 117, 1, 113, 113, 113, 1, 
	115, 118, 13, 13, 38, 38, 119, 38, 
	14, 38, 38, 115, 115, 115, 115, 115, 
	1, 120, 120, 120, 1, 115, 115, 115, 
	1, 121, 121, 121, 1, 119, 119, 119, 
	1, 115, 118, 13, 13, 38, 38, 13, 
	38, 14, 38, 38, 115, 115, 115, 115, 
	115, 1, 122, 122, 122, 1, 116, 116, 
	116, 1, 113, 114, 11, 11, 115, 11, 
	11, 14, 38, 38, 113, 113, 113, 113, 
	113, 1, 11, 12, 13, 11, 11, 11, 
	14, 11, 11, 11, 123, 93, 93, 1, 
	11, 12, 88, 124, 13, 11, 11, 11, 
	14, 11, 11, 11, 125, 91, 91, 1, 
	11, 12, 13, 11, 11, 11, 14, 11, 
	11, 11, 126, 93, 93, 1, 11, 12, 
	88, 127, 13, 11, 11, 11, 14, 11, 
	11, 11, 128, 91, 91, 1, 11, 12, 
	13, 11, 11, 11, 14, 11, 11, 11, 
	129, 93, 93, 1, 11, 12, 88, 127, 
	13, 11, 11, 11, 14, 11, 11, 11, 
	130, 91, 91, 1, 11, 12, 88, 127, 
	13, 11, 11, 11, 14, 11, 11, 11, 
	91, 91, 91, 1, 11, 12, 88, 124, 
	13, 11, 11, 11, 14, 11, 11, 11, 
	131, 91, 91, 1, 11, 12, 88, 124, 
	13, 11, 11, 11, 14, 11, 11, 11, 
	91, 91, 91, 1, 11, 12, 88, 89, 
	13, 11, 11, 11, 14, 11, 11, 11, 
	132, 91, 91, 1, 11, 12, 88, 89, 
	13, 11, 11, 11, 14, 11, 11, 11, 
	91, 91, 91, 1, 4, 1, 26, 133, 
	134, 135, 136, 25, 25, 25, 1, 134, 
	135, 136, 23, 25, 25, 1, 135, 136, 
	137, 1, 135, 136, 138, 1, 135, 136, 
	139, 1, 135, 136, 140, 1, 135, 136, 
	1, 31, 32, 141, 142, 143, 31, 31, 
	31, 31, 31, 31, 31, 1, 33, 34, 
	141, 143, 33, 33, 33, 33, 33, 33, 
	33, 1, 40, 144, 145, 40, 40, 40, 
	40, 40, 40, 40, 40, 1, 20, 24, 
	134, 135, 136, 146, 23, 23, 1, 20, 
	24, 134, 135, 136, 147, 23, 23, 1, 
	20, 24, 134, 135, 136, 23, 23, 23, 
	1, 134, 135, 136, 1, 11, 12, 94, 
	148, 149, 150, 11, 151, 14, 11, 11, 
	11, 93, 93, 93, 1, 11, 12, 149, 
	150, 11, 151, 14, 11, 11, 11, 91, 
	93, 93, 1, 13, 16, 135, 13, 136, 
	14, 13, 13, 13, 152, 13, 13, 1, 
	13, 16, 135, 13, 136, 14, 13, 13, 
	13, 153, 13, 13, 1, 13, 16, 135, 
	13, 136, 14, 13, 13, 13, 154, 13, 
	13, 1, 13, 16, 135, 13, 136, 14, 
	13, 13, 13, 155, 13, 13, 1, 13, 
	16, 135, 13, 136, 14, 13, 13, 13, 
	13, 13, 13, 1, 100, 106, 11, 102, 
	156, 157, 158, 14, 31, 31, 100, 100, 
	100, 100, 100, 1, 102, 159, 13, 31, 
	31, 141, 160, 143, 14, 31, 31, 102, 
	102, 102, 102, 102, 1, 103, 104, 13, 
	33, 33, 141, 13, 143, 14, 33, 33, 
	103, 103, 103, 103, 103, 1, 107, 108, 
	11, 103, 156, 11, 158, 14, 33, 33, 
	107, 107, 107, 107, 107, 1, 119, 161, 
	162, 13, 40, 40, 13, 40, 14, 40, 
	40, 119, 119, 119, 119, 119, 1, 116, 
	163, 164, 11, 119, 11, 11, 14, 40, 
	40, 116, 116, 116, 116, 116, 1, 11, 
	12, 88, 92, 149, 150, 11, 151, 14, 
	11, 11, 11, 165, 91, 91, 1, 11, 
	12, 88, 92, 149, 150, 11, 151, 14, 
	11, 11, 11, 166, 91, 91, 1, 11, 
	12, 88, 92, 149, 150, 11, 151, 14, 
	11, 11, 11, 91, 91, 91, 1, 
}

var _uri_trans_targs []byte = []byte{
	2, 0, 3, 4, 5, 121, 6, 7, 
	79, 136, 44, 6, 7, 9, 12, 8, 
	10, 11, 13, 122, 14, 33, 42, 15, 
	16, 122, 17, 124, 129, 20, 21, 129, 
	20, 130, 24, 25, 27, 28, 27, 28, 
	131, 29, 31, 34, 35, 40, 36, 37, 
	38, 132, 39, 41, 43, 45, 78, 46, 
	49, 135, 47, 48, 50, 65, 51, 63, 
	52, 53, 61, 54, 55, 59, 56, 57, 
	58, 60, 62, 64, 66, 74, 67, 70, 
	68, 69, 71, 72, 73, 75, 76, 77, 
	80, 110, 119, 81, 82, 136, 83, 138, 
	143, 86, 144, 87, 143, 89, 144, 145, 
	91, 92, 86, 146, 95, 96, 98, 99, 
	101, 98, 99, 101, 148, 100, 102, 147, 
	103, 105, 108, 111, 112, 117, 113, 114, 
	115, 149, 116, 118, 120, 123, 18, 19, 
	26, 125, 126, 127, 128, 22, 23, 26, 
	30, 32, 133, 134, 137, 84, 85, 97, 
	139, 140, 141, 142, 93, 94, 97, 88, 
	90, 104, 106, 107, 109, 150, 151, 
}

var _uri_trans_actions []byte = []byte{
	0, 0, 0, 0, 5, 0, 3, 3, 
	15, 15, 1, 0, 0, 0, 7, 0, 
	0, 0, 1, 1, 0, 0, 0, 0, 
	0, 0, 0, 0, 1, 1, 0, 0, 
	0, 0, 0, 0, 1, 1, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	1, 1, 1, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 1, 1, 
	1, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 9, 
	9, 0, 0, 0, 0, 0, 0, 11, 
	0, 0, 0, 0, 0, 0, 9, 9, 
	0, 0, 0, 0, 0, 0, 11, 0, 
	0, 0, 0, 0, 0, 0, 0, 
}

var _uri_eof_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 0, 0, 0, 0, 0, 0, 
	0, 0, 9, 9, 9, 9, 9, 9, 
	9, 11, 11, 13, 9, 9, 9, 9, 
	9, 9, 9, 9, 9, 9, 9, 11, 
	11, 11, 11, 13, 13, 9, 9, 9, 
}

const uri_start int = 1
const uri_first_final int = 122
const uri_error int = 0

const uri_en_main int = 1


//line parse_uri.rl:14
	uri := &URI{}
	m   := 0 // marker
	m1  := 0 // additional marker
	cs  := 0 // current state
	p   := 0 // data pointer
	pe  := len(data) // data end pointer
	eof := len(data)

	
//line parse_uri.rl:39


	
//line parse_uri.go:547
	{
	cs = uri_start
	}

//line parse_uri.rl:42
	
//line parse_uri.go:554
	{
	var _klen int
	var _trans int
	var _acts int
	var _nacts uint
	var _keys int
	if p == pe {
		goto _test_eof
	}
	if cs == 0 {
		goto _out
	}
_resume:
	_keys = int(_uri_key_offsets[cs])
	_trans = int(_uri_index_offsets[cs])

	_klen = int(_uri_single_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + _klen - 1)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + ((_upper - _lower) >> 1)
			switch {
			case data[p] < _uri_trans_keys[_mid]:
				_upper = _mid - 1
			case data[p] > _uri_trans_keys[_mid]:
				_lower = _mid + 1
			default:
				_trans += int(_mid - int(_keys))
				goto _match
			}
		}
		_keys += _klen
		_trans += _klen
	}

	_klen = int(_uri_range_lengths[cs])
	if _klen > 0 {
		_lower := int(_keys)
		var _mid int
		_upper := int(_keys + (_klen << 1) - 2)
		for {
			if _upper < _lower {
				break
			}

			_mid = _lower + (((_upper - _lower) >> 1) & ^1)
			switch {
			case data[p] < _uri_trans_keys[_mid]:
				_upper = _mid - 2
			case data[p] > _uri_trans_keys[_mid + 1]:
				_lower = _mid + 2
			default:
				_trans += int((_mid - int(_keys)) >> 1)
				goto _match
			}
		}
		_trans += _klen
	}

_match:
	_trans = int(_uri_indicies[_trans])
	cs = int(_uri_trans_targs[_trans])

	if _uri_trans_actions[_trans] == 0 {
		goto _again
	}

	_acts = int(_uri_trans_actions[_trans])
	_nacts = uint(_uri_actions[_acts]); _acts++
	for ; _nacts > 0; _nacts-- {
		_acts++
		switch _uri_actions[_acts-1] {
		case 0:
//line parse_uri.rl:23
 m = p 
		case 1:
//line parse_uri.rl:24
 m1 = p 
		case 2:
//line parse_uri.rl:25
 uri.Scheme   = data[:p] 
		case 3:
//line parse_uri.rl:26
 uri.Userinfo = data[m1:p] 
		case 4:
//line parse_uri.rl:27
 uri.Hostport = data[m:p] 
		case 5:
//line parse_uri.rl:28
 uri.Params   = Params(data[m:p]).setup() 
//line parse_uri.go:651
		}
	}

_again:
	if cs == 0 {
		goto _out
	}
	p++
	if p != pe {
		goto _resume
	}
	_test_eof: {}
	if p == eof {
		__acts := _uri_eof_actions[cs]
		__nacts := uint(_uri_actions[__acts]); __acts++
		for ; __nacts > 0; __nacts-- {
			__acts++
			switch _uri_actions[__acts-1] {
			case 4:
//line parse_uri.rl:27
 uri.Hostport = data[m:p] 
			case 5:
//line parse_uri.rl:28
 uri.Params   = Params(data[m:p]).setup() 
			case 6:
//line parse_uri.rl:29
 uri.Headers  = data[m:p] 
//line parse_uri.go:679
			}
		}
	}

	_out: {}
	}

//line parse_uri.rl:43

	if cs >= uri_first_final {
		return uri, nil
	}

	if p == pe {
		return nil, fmt.Errorf("%w: unexpected eof: %q", ErrURIParse, data)
	}

	return nil, fmt.Errorf("%w: error in uri at pos %d: %q>>%q<<", ErrURIParse, p, data[:p],data[p:])
}
