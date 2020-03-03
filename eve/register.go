// Copyright 2019 Michal Derkacz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eve

type Register uint32

// Register list, offset from RAM_REG: FT81 FT80
const (
	REG_ID                Register = 0x0000_0000
	REG_FRAMES            Register = 0x0004_0004
	REG_CLOCK             Register = 0x0008_0008
	REG_FREQUENCY         Register = 0x000C_000C
	REG_RENDERMODE        Register = 0x0010_0010
	REG_SNAPY             Register = 0x0014_0014
	REG_SNAPSHOT          Register = 0x0018_0018
	REG_SNAPFORMAT        Register = 0x001C_FFFF
	REG_CPURESET          Register = 0x0020_001C
	REG_TAP_CRC           Register = 0x0024_0020
	REG_TAP_MASK          Register = 0x0028_0024
	REG_HCYCLE            Register = 0x002C_0028
	REG_HOFFSET           Register = 0x0030_002C
	REG_HSIZE             Register = 0x0034_0030
	REG_HSYNC0            Register = 0x0038_0034
	REG_HSYNC1            Register = 0x003C_0038
	REG_VCYCLE            Register = 0x0040_003C
	REG_VOFFSET           Register = 0x0044_0040
	REG_VSIZE             Register = 0x0048_0044
	REG_VSYNC0            Register = 0x004C_0048
	REG_VSYNC1            Register = 0x0050_004C
	REG_DLSWAP            Register = 0x0054_0050
	REG_ROTATE            Register = 0x0058_0054
	REG_OUTBITS           Register = 0x005C_0058
	REG_DITHER            Register = 0x0060_005C
	REG_SWIZZLE           Register = 0x0064_0060
	REG_CSPREAD           Register = 0x0068_0064
	REG_PCLK_POL          Register = 0x006C_0068
	REG_PCLK              Register = 0x0070_006C
	REG_TAG_X             Register = 0x0074_0070
	REG_TAG_Y             Register = 0x0078_0074
	REG_TAG               Register = 0x007C_0078
	REG_VOL_PB            Register = 0x0080_007C
	REG_VOL_SOUND         Register = 0x0084_0080
	REG_SOUND             Register = 0x0088_0084
	REG_PLAY              Register = 0x008C_0088
	REG_GPIO_DIR          Register = 0x0090_008C
	REG_GPIO              Register = 0x0094_0090
	REG_GPIOX_DIR         Register = 0x0098_FFFF
	REG_GPIOX             Register = 0x009C_FFFF
	REG_INT_FLAGS         Register = 0x00A8_0098
	REG_INT_EN            Register = 0x00AC_009C
	REG_INT_MASK          Register = 0x00B0_00A0
	REG_PLAYBACK_START    Register = 0x00B4_00A4
	REG_PLAYBACK_LENGTH   Register = 0x00B8_00A8
	REG_PLAYBACK_READPTR  Register = 0x00BC_00AC
	REG_PLAYBACK_FREQ     Register = 0x00C0_00B0
	REG_PLAYBACK_FORMAT   Register = 0x00C4_00B4
	REG_PLAYBACK_LOOP     Register = 0x00C8_00B8
	REG_PLAYBACK_PLAY     Register = 0x00CC_00BC
	REG_PWM_HZ            Register = 0x00D0_00C0
	REG_PWM_DUTY          Register = 0x00D4_00C4
	REG_MACRO_0           Register = 0x00D8_00C8
	REG_MACRO_1           Register = 0x00DC_00CC
	REG_CMD_READ          Register = 0x00F8_00E4
	REG_CMD_WRITE         Register = 0x00FC_00E8
	REG_CMD_DL            Register = 0x0100_00EC
	REG_TOUCH_MODE        Register = 0x0104_00F0
	REG_CTOUCH_EXTENDED   Register = 0x0108_00F4
	REG_TOUCH_ADC_MODE    Register = 0x0108_00F4
	REG_TOUCH_CHARGE      Register = 0x010C_00F8
	REG_CTOUCH_REG        Register = 0xFFFF_00F8
	REG_TOUCH_SETTLE      Register = 0x0110_00FC
	REG_TOUCH_OVERSAMPLE  Register = 0x0114_0100
	REG_TOUCH_RZTHRESH    Register = 0x0118_0104
	REG_TOUCH_RAW_XY      Register = 0x011C_0108
	REG_CTOUCH_TOUCH1_XY  Register = 0x011C_0108
	REG_CTOUCH_TOUCH4_Y   Register = 0x0120_010C
	REG_TOUCH_RZ          Register = 0x0120_010C
	REG_CTOUCH_TOUCH0_XY  Register = 0x0124_0110
	REG_TOUCH_SCREEN_XY   Register = 0x0124_0110
	REG_TOUCH_TAG_XY      Register = 0x0128_0114
	REG_TOUCH_TAG         Register = 0x012C_0118
	REG_TOUCH_TAG1_XY     Register = 0x0130_FFFF
	REG_TOUCH_TAG1        Register = 0x0134_FFFF
	REG_TOUCH_TAG2_XY     Register = 0x0138_FFFF
	REG_TOUCH_TAG2        Register = 0x013C_FFFF
	REG_TOUCH_TAG3_XY     Register = 0x0140_FFFF
	REG_TOUCH_TAG3        Register = 0x0144_FFFF
	REG_TOUCH_TAG4_XY     Register = 0x0148_FFFF
	REG_TOUCH_TAG4        Register = 0x014C_FFFF
	REG_TOUCH_TRANSFORM_A Register = 0x0150_011C
	REG_TOUCH_TRANSFORM_B Register = 0x0154_0120
	REG_TOUCH_TRANSFORM_C Register = 0x0158_0124
	REG_TOUCH_TRANSFORM_D Register = 0x015C_0128
	REG_TOUCH_TRANSFORM_E Register = 0x0160_012C
	REG_TOUCH_TRANSFORM_F Register = 0x0164_0130
	REG_TOUCH_CONFIG      Register = 0x0168_FFFF
	REG_CTOUCH_TOUCH4_X   Register = 0x016C_0138
	REG_BIST_EN           Register = 0x0174_FFFF
	REG_TRIM              Register = 0x0180_016C
	REG_ANA_COMP          Register = 0x0184_FFFF
	REG_SPI_WIDTH         Register = 0x0188_FFFF
	REG_CTOUCH_TOUCH2_XY  Register = 0x018C_0174
	REG_TOUCH_DIRECT_XY   Register = 0x018C_0174
	REG_TOUCH_DIRECT_Z1Z2 Register = 0x0190_0178
	REG_CTOUCH_TOUCH3_XY  Register = 0x0190_0178
	REG_DATESTAMP         Register = 0x0564_FFFF
	REG_CMDB_SPACE        Register = 0x0574_FFFF
	REG_CMDB_WRITE        Register = 0x0578_FFFF
	REG_TRACKER           Register = 0x7000_6C00
	REG_TRACKER_1         Register = 0x7004_FFFF
	REG_TRACKER_2         Register = 0x7008_FFFF
	REG_TRACKER_3         Register = 0x700C_FFFF
	REG_TRACKER_4         Register = 0x7010_FFFF
	REG_MEDIAFIFO_READ    Register = 0x7014_FFFF
	REG_MEDIAFIFO_WRITE   Register = 0x7018_FFFF
)

// REG_DLSWAP values.
const (
	DLSWAP_DONE  = 0
	DLSWAP_LINE  = 1
	DLSWAP_FRAME = 2
)
