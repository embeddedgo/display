package invaders

const assets = "" +
	"\x78\x9C\xED\x91\xB1\x4A\xC4\x40\x10\x86\x27\x9B\x68\xAE\x38\xC8" +
	"\x26\x95\xA2\xB8\x6B\x61\x79\x70\x36\x8A\x8D\x49\x23\x1C\x58\x29" +
	"\x3E\x80\x95\xDD\x89\x6F\x90\xD5\x27\xD0\xDE\x22\x60\x63\x69\x65" +
	"\x7B\x5B\x58\xF8\x18\x0B\x8A\xA5\x44\x41\x88\x70\xDC\xB8\x93\x5C" +
	"\xCE\x13\x41\x04\x11\xF5\xF0\x1F\x98\xE4\xCF\xF7\xEF\x64\x37\x69" +
	"\xC2\x02\xC4\xD0\x55\xA9\x8A\xA0\x05\xDB\xAA\x39\xF2\x2D\x48\xD4" +
	"\x22\x48\x75\x22\x7B\xD9\xBD\x41\x93\x6A\xA9\x92\x44\xAA\x10\x62" +
	"\xB5\xAF\xD1\xF4\xB2\x2C\x9B\x83\x00\x52\x83\xF9\xDD\x01\xE6\xB3" +
	"\x6A\x45\xEB\xF6\xB8\xDF\xDB\x6D\xEB\x86\xDA\xE0\x52\xF2\x44\xCF" +
	"\xF0\x64\x49\x76\x8E\x19\xF8\xB6\xD2\x1C\x8B\xAA\xB8\x77\x34\xCB" +
	"\x35\x63\x9B\xCB\x8B\x57\x22\x5F\x7D\x64\x00\x9C\x65\x11\x38\x17" +
	"\x1B\xE7\xE2\x36\x7E\x62\x50\x55\xAD\x50\x7A\x49\xA0\x85\x11\x26" +
	"\xD0\x9E\x8A\x24\x03\xCF\x52\x07\xE8\xEA\x58\xEE\x81\x3B\x74\xAC" +
	"\xF4\x6C\xCC\xD1\x94\x71\xE7\x95\xBC\x2E\xBF\x9C\x5F\xDF\xD7\xEF" +
	"\x7C\x75\x95\xF7\xC7\x56\xC0\x9B\xF5\xD5\x13\xEA\xD3\x76\x0F\x55" +
	"\x7F\xCF\x5D\xFB\x9C\x0D\x7B\xB5\x1F\x2F\x71\x32\x57\x4F\x69\xD7" +
	"\x4C\x69\x96\xD1\xCC\x00\xB5\xED\xE9\x80\xE2\x98\xDB\x2E\x50\xD5" +
	"\xE7\x17\x68\x6C\x47\xA4\xAD\x60\x61\x7B\x5C\x86\x4B\xC5\x14\x76" +
	"\x90\xC2\x01\xF6\x69\x44\x19\x2E\x95\x52\xD8\x45\x0A\x0B\x1C\x94" +
	"\x23\xF2\x9A\x21\x85\x7D\xA4\x70\x4C\x93\xED\x88\x62\x92\xD8\xC0" +
	"\x25\x96\x0B\x62\x8A\x58\xDF\x1F\xB1\xBE\x43\xCC\x04\xC4\x20\xB5" +
	"\xAC\x70\x47\xAC\x00\x62\xDA\xB7\x6C\xDD\x52\x70\x72\xE7\x53\x6C" +
	"\xB2\x15\xF9\xA6\x21\x8A\xAD\xB8\xEF\x76\x2F\xED\x07\x0F\xD7\x4E" +
	"\x13\xDE\x08\x0F\x3B\x30\xCD\xE1\x06\xE0\xF9\x1A\x76\x6E\xE1\x41" +
	"\xC1\x19\x07\xAE\xEC\x7F\xF9\x3D\x32\xDF\xA0\x9F\x3E\xD3\xBF\xBE" +
	"\xA6\xF9\x3F\xA2\x8F\x76\xFA\xD1\xF9\x5E\x00\xA4\xC6\xB4\xF2"

const (
	spr16Addr   = 0
	spr16Width  = 16
	spr16Height = 200

	shieldsAddr   = 400
	shieldsWidth  = 224
	shieldsHeight = 24

	saucerAddr   = 1072
	saucerWidth  = 24
	saucerHeight = 16

	overlayAddr   = 1120
	overlayWidth  = 28
	overlayHeight = 32

	assetsEnd = 2016
)

const background_jpg = "" +
	"\xFF\xD8\xFF\xE0\x00\x10\x4A\x46\x49\x46\x00\x01\x01\x00\x00\x01" +
	"\x00\x01\x00\x00\xFF\xDB\x00\x43\x00\x0E\x0A\x0B\x0D\x0B\x09\x0E" +
	"\x0D\x0C\x0D\x10\x0F\x0E\x11\x16\x24\x17\x16\x14\x14\x16\x2C\x20" +
	"\x21\x1A\x24\x34\x2E\x37\x36\x33\x2E\x32\x32\x3A\x41\x53\x46\x3A" +
	"\x3D\x4E\x3E\x32\x32\x48\x62\x49\x4E\x56\x58\x5D\x5E\x5D\x38\x45" +
	"\x66\x6D\x65\x5A\x6C\x53\x5B\x5D\x59\xFF\xDB\x00\x43\x01\x0F\x10" +
	"\x10\x16\x13\x16\x2A\x17\x17\x2A\x59\x3B\x32\x3B\x59\x59\x59\x59" +
	"\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59" +
	"\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59" +
	"\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\x59\xFF\xC0" +
	"\x00\x11\x08\x01\x10\x00\xF8\x03\x01\x22\x00\x02\x11\x01\x03\x11" +
	"\x01\xFF\xC4\x00\x1F\x00\x00\x01\x05\x01\x01\x01\x01\x01\x01\x00" +
	"\x00\x00\x00\x00\x00\x00\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09" +
	"\x0A\x0B\xFF\xC4\x00\xB5\x10\x00\x02\x01\x03\x03\x02\x04\x03\x05" +
	"\x05\x04\x04\x00\x00\x01\x7D\x01\x02\x03\x00\x04\x11\x05\x12\x21" +
	"\x31\x41\x06\x13\x51\x61\x07\x22\x71\x14\x32\x81\x91\xA1\x08\x23" +
	"\x42\xB1\xC1\x15\x52\xD1\xF0\x24\x33\x62\x72\x82\x09\x0A\x16\x17" +
	"\x18\x19\x1A\x25\x26\x27\x28\x29\x2A\x34\x35\x36\x37\x38\x39\x3A" +
	"\x43\x44\x45\x46\x47\x48\x49\x4A\x53\x54\x55\x56\x57\x58\x59\x5A" +
	"\x63\x64\x65\x66\x67\x68\x69\x6A\x73\x74\x75\x76\x77\x78\x79\x7A" +
	"\x83\x84\x85\x86\x87\x88\x89\x8A\x92\x93\x94\x95\x96\x97\x98\x99" +
	"\x9A\xA2\xA3\xA4\xA5\xA6\xA7\xA8\xA9\xAA\xB2\xB3\xB4\xB5\xB6\xB7" +
	"\xB8\xB9\xBA\xC2\xC3\xC4\xC5\xC6\xC7\xC8\xC9\xCA\xD2\xD3\xD4\xD5" +
	"\xD6\xD7\xD8\xD9\xDA\xE1\xE2\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xF1" +
	"\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9\xFA\xFF\xC4\x00\x1F\x01\x00\x03" +
	"\x01\x01\x01\x01\x01\x01\x01\x01\x01\x00\x00\x00\x00\x00\x00\x01" +
	"\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\xFF\xC4\x00\xB5\x11\x00" +
	"\x02\x01\x02\x04\x04\x03\x04\x07\x05\x04\x04\x00\x01\x02\x77\x00" +
	"\x01\x02\x03\x11\x04\x05\x21\x31\x06\x12\x41\x51\x07\x61\x71\x13" +
	"\x22\x32\x81\x08\x14\x42\x91\xA1\xB1\xC1\x09\x23\x33\x52\xF0\x15" +
	"\x62\x72\xD1\x0A\x16\x24\x34\xE1\x25\xF1\x17\x18\x19\x1A\x26\x27" +
	"\x28\x29\x2A\x35\x36\x37\x38\x39\x3A\x43\x44\x45\x46\x47\x48\x49" +
	"\x4A\x53\x54\x55\x56\x57\x58\x59\x5A\x63\x64\x65\x66\x67\x68\x69" +
	"\x6A\x73\x74\x75\x76\x77\x78\x79\x7A\x82\x83\x84\x85\x86\x87\x88" +
	"\x89\x8A\x92\x93\x94\x95\x96\x97\x98\x99\x9A\xA2\xA3\xA4\xA5\xA6" +
	"\xA7\xA8\xA9\xAA\xB2\xB3\xB4\xB5\xB6\xB7\xB8\xB9\xBA\xC2\xC3\xC4" +
	"\xC5\xC6\xC7\xC8\xC9\xCA\xD2\xD3\xD4\xD5\xD6\xD7\xD8\xD9\xDA\xE2" +
	"\xE3\xE4\xE5\xE6\xE7\xE8\xE9\xEA\xF2\xF3\xF4\xF5\xF6\xF7\xF8\xF9" +
	"\xFA\xFF\xDA\x00\x0C\x03\x01\x00\x02\x11\x03\x11\x00\x3F\x00\xE1" +
	"\x28\xC5\x29\xF6\xA5\xDA\x72\x07\x5C\xFA\x56\xC4\x02\x26\xE3\xD3" +
	"\x35\x21\x85\x40\xC0\x19\x6A\x99\x57\xCA\x4E\x39\x26\x97\x25\x53" +
	"\x2F\xC0\xFD\x4D\x66\xD8\x15\xCC\x69\x18\xF9\x80\x66\xF4\xA8\x5B" +
	"\x04\xE7\x03\xF0\xA7\x33\x6E\x39\xA6\xD5\x24\x02\x60\x7A\x51\x81" +
	"\xE9\x4E\xC7\x14\x94\x58\x77\x13\x03\xD2\x8C\x66\x96\x8A\x76\x01" +
	"\x30\x3D\x29\x70\x3D\x29\x68\xA2\xC2\x0C\x7B\x51\x8A\x5A\x5A\x76" +
	"\x01\x31\xDF\x14\x52\xD2\xD3\xB0\x08\x79\xE7\xBD\x25\x2D\x06\x95" +
	"\x80\x43\xD7\x81\xC5\x26\x29\xC0\xE3\xA7\xA5\x25\x00\x26\x29\x68" +
	"\xA2\x80\x0A\x29\x69\x28\x00\xA4\xA5\xA2\x81\x8D\xA2\x96\x8A\x40" +
	"\x36\x8A\x5A\x43\x48\x04\x34\x51\x45\x48\xC9\xB0\x0E\x00\xCE\x6A" +
	"\xD4\x36\xFB\x7E\x67\xA8\xED\x50\x3C\x9C\xF4\x1C\xD5\x86\x95\x4C" +
	"\x9C\x9C\x28\xFD\x6A\xA4\xFA\x22\x05\x7D\xA8\x37\x37\x41\x54\xE4" +
	"\x90\xCA\xD8\xE7\x6F\xA0\xA9\x24\x6F\x3A\x4D\xAB\x9C\x54\x8B\x6A" +
	"\x01\xCE\x6A\x56\x9B\x8C\xAE\x90\x31\x3C\x8C\x54\xA2\xDC\x0E\xF9" +
	"\xA9\xCA\xAA\xF4\xA8\xCB\x1D\xD8\x06\x95\xDB\x02\x19\x23\x0A\x38" +
	"\x15\x11\x43\x8C\xE2\xAC\xC9\x90\xB9\xA2\x33\xBD\x08\x22\x9A\x62" +
	"\x29\x9A\x29\xCC\x30\x69\xB5\xA0\xD0\x52\xD1\x45\x00\x2D\x02\x8A" +
	"\x5A\x62\x01\x45\x14\xB4\xC0\x4A\x4C\x53\xB1\xC5\x25\x20\x12\x8A" +
	"\x28\xA4\x31\x28\xA5\x23\x9E\x28\xA6\x02\x51\x4B\x45\x02\x12\x8A" +
	"\x28\xA4\x31\x29\x68\xA2\x80\x13\x1C\x66\x90\xD3\xA9\x0D\x20\x1B" +
	"\x45\x06\x8A\x91\x96\x51\x8A\x03\x8E\xF4\xC3\xC9\xA9\x36\xFC\xB9" +
	"\xF4\xA4\x1B\x01\x6E\xE3\x1C\x55\x99\xA1\xD6\xA4\x09\x39\xAB\xE7" +
	"\xA5\x66\x03\xCE\x6A\x71\x74\xC1\x70\x79\xA8\x94\x6E\xCA\xB9\x64" +
	"\xE3\x15\x4F\x39\x92\x91\xA6\x63\xDE\x91\x17\x9E\x45\x0A\x36\x13" +
	"\x65\x92\xA0\xAD\x57\x0D\xB1\xCD\x4E\xD2\x80\xB8\xAA\xAF\xD4\xD1" +
	"\x14\x0D\x8C\x6E\x4D\x25\x2D\x15\x63\x12\x96\x8A\x29\x80\x52\xD1" +
	"\x4B\x4C\x02\x8A\x06\x33\xCD\x07\xAD\x00\x07\x9A\x4A\x5A\x4A\x00" +
	"\x4A\x5C\x71\x46\x29\x71\x48\x04\xA5\x23\x14\xA7\x1D\xB9\xA6\xD0" +
	"\x02\x51\x4B\x4A\x79\xF4\xF4\xA0\x06\xE2\x8A\x5A\x28\x01\x28\xA5" +
	"\xC5\x18\xA0\x04\xA4\x34\xEA\x69\xA1\x8C\x69\xA2\x83\x45\x40\xCB" +
	"\xCA\xBB\xA1\x6C\x73\x83\xCD\x47\x24\x4C\x84\x12\x38\x3D\x08\xA7" +
	"\xC6\xFB\x33\xEF\x56\xA3\x25\xA3\xC1\x5D\xCB\xEF\x4D\xDD\x19\x23" +
	"\x37\x1C\xD2\xED\x20\xE2\xAD\xCA\x89\xD5\x46\x0D\x40\x41\xCF\x34" +
	"\xD3\xB8\x36\x33\x6D\x48\x58\x81\x82\x69\x17\x00\xF2\x33\x43\xB6" +
	"\x7B\x00\x07\xA5\x16\x0B\x8D\x2C\x73\x9E\xF4\xC3\xC9\xC9\xA5\xEA" +
	"\x79\xA4\xA7\x62\x90\x86\x93\x14\xEA\x31\x4C\x63\x71\x4B\x83\x8C" +
	"\xE0\xE3\xD6\x96\x9F\xBD\x8C\x41\x09\xF9\x41\xCE\x29\x01\x1D\x28" +
	"\xE0\x7D\x68\xA5\xC5\x31\x09\x4B\x46\x29\x68\x01\x29\x29\xD8\xA4" +
	"\xC5\x30\x01\x47\x5A\x51\xC1\x04\x75\x14\x75\x3C\xD2\x0B\x89\x49" +
	"\x4E\xA4\xC5\x00\x25\x28\xE0\x83\x45\x18\xA2\xC0\x27\x5A\x29\x71" +
	"\x45\x00\x25\x14\xEE\xDD\x3F\x1A\x4C\x50\x02\x1E\x80\x76\xA6\x9E" +
	"\x98\xA7\xE2\x9A\x45\x16\x18\xC3\x45\x29\xA2\xA1\x8C\xB0\x0D\x59" +
	"\x82\x7D\x83\x07\xA5\x56\xC6\x3B\xD1\xCD\x5B\x57\x32\xB5\xB6\x2E" +
	"\x3D\xCA\x63\xE5\x5A\xAA\xEE\x59\xB2\x69\xB9\xA5\x50\x49\x3E\xC3" +
	"\x34\x94\x52\x0D\x58\x8C\x47\x6A\x61\xA5\xA0\x8F\x4A\xA2\x92\x1B" +
	"\x45\x2D\x18\xA0\x62\x51\x4B\x8A\x31\x40\x09\x8A\x31\xCD\x2E\x29" +
	"\x71\x40\xAE\x25\x18\xA7\x62\x8C\x50\x17\x13\x14\x53\xB1\x46\x28" +
	"\x15\xC4\xC5\x18\xA5\xC7\x34\x50\x31\xA4\x51\x8A\x76\x28\xC5\x00" +
	"\x37\x14\x62\x9D\x45\x01\x71\xB8\xA3\x14\xE2\x29\x31\x40\x09\x8A" +
	"\x31\x4B\xDF\x8A\x5A\x60\x21\xF4\x19\xC7\xBD\x00\x0C\xF3\x9C\x7B" +
	"\x53\x94\x64\xE0\x90\x07\xA9\xA4\x22\x90\x0C\xC5\x21\xA7\xE2\x9A" +
	"\x45\x03\x4C\x8C\x8A\x29\x58\x51\x52\x51\x35\x2D\x28\xEB\x4F\x11" +
	"\x97\x5C\xA7\x2C\x49\xCA\x0E\xA0\x7A\xD5\x99\x91\x51\x8A\x76\x29" +
	"\x28\x10\x98\xA4\xC5\x3A\x8C\x50\x31\xB8\xA0\xF3\x8F\x6A\x76\x29" +
	"\x31\x40\x09\x8A\x31\x4E\xC7\x7A\x4A\x00\x4C\x52\xE2\x96\x8A\x00" +
	"\x31\x46\x29\x68\xA7\x60\x12\x8A\x76\x28\xA2\xC0\x37\x14\xBD\xA9" +
	"\x71\x46\x28\xB0\x86\xE2\x8A\x76\x28\xC5\x03\x12\x8A\x31\x4B\x45" +
	"\x80\x41\xEF\x48\x45\x38\x51\x83\xE9\x48\x06\xD1\x8E\x69\x71\x46" +
	"\x28\x01\x31\x45\x2D\x14\x00\x98\xE9\xCE\x7F\xA5\x34\x8A\x7E\x29" +
	"\x28\x02\x13\x45\x39\x85\x15\x2C\xA2\x6A\x74\x6E\xD1\xB6\xE4\x62" +
	"\xA7\x18\xC8\xA4\xA5\xAB\x20\x6D\x18\xA5\xC5\x18\xA0\x04\x34\x98" +
	"\xA7\x51\x40\x0D\xC5\x18\xA7\x52\x62\x80\x1F\x24\x5B\x23\x8D\xF7" +
	"\xA9\xDE\x3A\x0E\xD5\x1E\x29\xD8\xA3\x14\x58\x2E\x37\x14\x01\x4E" +
	"\x23\x06\x8C\x50\x00\x05\x18\xA5\xA5\xA6\x21\x31\x46\x29\xD8\xA5" +
	"\xC5\x01\x71\xB8\xA4\xC5\x3C\xE4\xE3\xDA\x82\x3A\xF6\xA0\x57\x1B" +
	"\x8A\x42\x30\x78\xA7\xE2\x90\x8E\x39\xA0\x63\x31\x46\x29\xD8\xA3" +
	"\x14\x0C\x4C\x52\x62\x9D\x8A\x28\x10\xDC\x52\xB0\x19\xF9\x73\x8F" +
	"\x7A\x5C\x50\x45\x16\x18\x84\x2E\xC5\xC6\x77\x7F\x17\xA5\x37\x14" +
	"\xEC\x51\x8A\x40\x34\xFD\x31\x49\x4F\xED\x8F\xC6\x9A\x45\x00\x46" +
	"\xE2\x8A\x71\x19\xA2\xA5\x94\x99\x25\x2D\x29\x1E\xF4\x55\x90\x25" +
	"\x25\x3B\x14\x62\x80\x1B\x4A\x47\xA5\x2D\x18\xA0\x06\xE2\x97\x14" +
	"\xB8\xA5\xC5\x02\xB8\xDC\x51\x8A\x79\x00\x9E\x06\x07\xA5\x18\xA0" +
	"\x2E\x33\x14\xBB\x70\x07\x1D\x69\xE5\x40\x40\x43\x64\x9E\xA3\x1D" +
	"\x28\x19\x14\x00\xCC\x53\x80\x3B\x7D\xA9\x71\x46\x28\x15\xC0\x0A" +
	"\x5C\x50\x05\x3B\x0C\x00\x24\x1C\x1E\x87\xD6\x98\x86\x81\x46\x29" +
	"\xDC\x81\xC8\x34\xF5\x71\x8C\x32\x86\x18\x38\x1D\x30\x7D\x68\x02" +
	"\x2C\x51\x8A\x93\x02\x90\x8A\x07\x72\x22\xB4\x98\xA9\x88\xA4\xC5" +
	"\x01\x72\x32\x07\x6C\xD1\x8A\x7E\x29\x31\xCD\x20\x19\x8A\x31\x4E" +
	"\xC5\x3F\xCA\x06\x03\x26\xF5\xC8\x38\xDB\xDF\xEB\x45\xCA\x21\xC5" +
	"\x18\xA7\x62\x83\x9C\x63\xB5\x00\x33\x14\x86\x9E\x45\x34\x8A\x00" +
	"\x6E\xC2\x51\x9B\x23\x0B\xD7\x9A\x28\x22\x8A\x9B\x0C\x90\x0A\x5E" +
	"\x87\x8A\x52\x30\x78\xE9\x45\x58\x84\xC5\x18\xA7\x52\xD0\x21\x98" +
	"\xA3\x14\xEC\x52\xE2\x81\x5C\x6E\x29\x40\x1D\xE9\xE0\x2E\xD3\x90" +
	"\x77\x76\xF4\xA4\xC5\x02\x1B\x8A\x5D\xB4\xF0\x29\xC0\x50\x4B\x63" +
	"\x10\x2E\xE1\xBF\x3B\x7D\xA9\x31\x52\x6D\xA3\x14\x05\xC6\x01\x46" +
	"\x2A\x4D\xB4\xBB\x68\x13\x64\x7B\x7B\xD1\xCE\x00\xC9\xC0\xE9\xED" +
	"\x52\xE0\xE3\x1D\x87\x6A\x4D\xB4\x0E\xE3\x0E\x49\x24\x9C\x93\x42" +
	"\xFC\xA7\xA6\x78\xA7\xE2\x8D\xB4\x58\x77\x23\xC5\x2D\x3C\x81\xDA" +
	"\x8C\x71\x40\x11\xD2\x10\x6A\x4C\x50\x57\x14\x01\x16\x28\xC5\x49" +
	"\x8A\x42\x28\x0B\x91\xE2\x8C\x53\xC8\xA4\xC5\x03\xB8\xDC\x52\x11" +
	"\x93\xEF\x4F\xC5\x26\x28\x1D\xC6\xB6\x09\xC8\x5D\xA3\xD2\x98\x45" +
	"\x49\x8A\x42\x29\x0E\xE4\x6C\xC4\xA8\x52\x78\x1D\x28\xA5\x22\x8A" +
	"\x2C\x3B\x8F\xC5\x28\x1C\x52\xE2\x97\x14\xC4\x20\x14\xB8\xF6\xA7" +
	"\x01\x4E\x0B\xC7\x4A\x08\x6C\x8F\x14\xF5\x42\xCD\x85\x19\x34\xA1" +
	"\x79\xA7\x6D\xA0\x9B\x91\x85\xA5\xC5\x48\x16\x97\x6D\x02\xE6\x23" +
	"\x03\x34\xF0\xB4\xF0\xB4\xE0\xB4\x12\xE4\x47\xB6\x9C\xCA\xB8\x5D" +
	"\xB9\xCE\x39\xCF\xAD\x3F\x6D\x2E\xDA\x09\xB9\x16\xDA\x50\xB5\x2E" +
	"\xDA\x5D\xB4\x05\xC8\x76\xD1\xB6\xA6\xDB\x46\xDA\x06\xA4\x43\xB6" +
	"\x8D\xB5\x36\xDA\x00\xE0\x8C\x0E\x7B\xFA\x53\x2B\x98\x87\x6D\x26" +
	"\xDA\x9F\x6D\x26\xDA\x06\xA4\x41\x8A\x08\x24\x92\x4E\x49\xA9\x76" +
	"\xE7\xA5\x21\x5A\x02\xE4\x44\x53\x71\x53\x62\x9A\x45\x03\xB9\x1E" +
	"\x29\xB8\xA9\x4A\xFA\x52\x11\x48\x64\x78\xA6\xE2\xA4\x22\x8D\xB9" +
	"\x1D\x40\xC5\x03\xB9\x1E\x29\x31\x52\x62\x90\x81\x48\x13\x22\x72" +
	"\x58\x01\x81\xC7\xA0\xA2\x95\x81\xA2\x82\xAE\x3C\x0E\x29\x40\xA5" +
	"\x02\x9E\x17\xD6\x82\x1B\x10\x2F\x7A\x70\x14\xA0\x71\x4A\x05\x04" +
	"\x36\x20\x14\xB8\xAB\x37\x36\xDF\x67\x28\x37\xAB\xEF\x5D\xDF\x2F" +
	"\x6A\x8B\x14\x27\x7D\x49\x77\x5A\x0C\x0B\x4E\x02\x9E\x05\x2E\x29" +
	"\x92\xD8\xDC\x52\xED\xA7\xE2\xA5\xD9\xE5\x3B\x2B\xA8\x73\x8E\x30" +
	"\xDC\x03\xEB\x40\x10\x05\xA5\xC5\x48\x17\x8E\x94\xBB\x68\x11\x18" +
	"\x5A\x50\xB5\x26\xDF\x6A\x5D\xB4\xC5\x72\x3D\xB8\xA3\x6D\x4B\xB6" +
	"\x8D\xB4\x0E\xE4\x5B\x68\xDB\x52\xE2\x95\x62\x66\xDD\xB4\x67\x68" +
	"\xC9\xF6\x14\x05\xC8\x36\xD1\xB6\xA6\xDB\x49\xB6\x80\xB9\x0E\xDA" +
	"\x6E\xDA\x9F\x6D\x21\x5A\x07\x72\x02\xB4\xD2\xB5\x39\x5A\x69\x5A" +
	"\x0A\x52\x20\xDB\x4D\x2B\x53\x95\xA6\xED\xA4\x55\xC8\x48\x1D\xB3" +
	"\x4D\x22\xA5\x22\x9A\x45\x03\xB8\xCC\x52\x30\x1B\x47\xAF\x7A\x7E" +
	"\x29\x18\x74\xA4\x08\x88\xA1\xE3\x20\xFA\x8A\x2A\xD1\x88\x32\x82" +
	"\x87\xE6\xC1\x25\x3F\xBA\x28\xA5\x70\x77\x21\x03\x15\x22\x8A\x58" +
	"\xCE\xD7\x0D\xB4\x36\x3B\x1E\x86\x94\x0A\x09\x60\x05\x28\x14\xE5" +
	"\x1E\xF8\xA7\x01\x4C\x91\xA1\x69\xDB\x70\x79\xC8\xA7\x8C\xED\xC7" +
	"\x6C\xE6\x9C\x06\x7A\xD0\x2B\x8C\x0B\x4F\xF2\x88\x45\x6C\x8C\x37" +
	"\x41\x9A\x70\x14\xC9\x65\x8A\x05\xCC\xAE\x17\xEB\xD4\xD3\xD8\x4B" +
	"\x5D\x10\xE0\xBC\x52\x85\xAC\xE9\x35\x52\x41\xFB\x34\x25\xBF\xDA" +
	"\x73\x81\xF9\x55\x19\xAF\x2E\x1C\x9F\x3E\xE0\xAA\xFF\x00\x75\x38" +
	"\x15\x8C\xAB\xC1\x6D\xA9\xD5\x0C\x25\x59\x6A\xF4\x46\xF4\xAF\x1C" +
	"\x0B\xBA\x59\x15\x07\xB9\xAA\xA7\x54\xB4\x07\xE5\x76\x73\xFE\xCA" +
	"\xD6\x14\x61\x58\x63\x6B\x49\xCF\x56\xAB\x00\x2E\xDC\x1C\x73\xD8" +
	"\x56\x12\xC5\x35\xB2\x3A\xE9\xE5\xC9\xEB\x26\x5D\x5D\x61\xD9\xCE" +
	"\xDB\x7F\x94\x1E\x09\x6A\x56\xD5\x66\x62\x76\xC2\x83\xEA\x49\xAA" +
	"\x58\x54\x5E\x70\xA2\xA3\x89\x94\xB7\x12\x17\xE7\xD2\xB2\x78\x8A" +
	"\x8F\xA9\xD2\xB0\x54\x56\xE8\xB0\xF7\xF7\x65\xF2\x66\xD8\x3D\x06" +
	"\x05\x1F\x6F\xB9\x6C\x01\x72\x01\xF6\x51\x55\xA4\x86\x21\x28\xDD" +
	"\x90\x5B\xF5\xA9\x56\x14\x5F\xBB\x8A\x97\x56\x5D\xCB\x58\x6A\x7F" +
	"\xCA\x8B\x2B\x7F\x74\x80\x9D\xEB\x20\xF4\x65\xA9\x57\x52\xB8\xEA" +
	"\x52\x3F\xA7\x35\x9F\x33\x88\xD7\x0C\x4F\x3E\x82\x92\x19\x15\x80" +
	"\x0A\xAF\x81\xDE\x85\x56\xA5\xB7\x13\xC3\x51\x72\xB7\x29\xA6\x9A" +
	"\xA9\xDF\xF3\x46\x85\x7B\xED\x6E\x45\x5A\x5B\xEB\x67\x03\xF7\x98" +
	"\xCF\xA8\xAC\x7C\x0E\x46\x06\x7D\x2A\x37\x6F\x28\x65\x63\x27\x3D" +
	"\x71\xDA\xAE\x38\x99\xA3\x29\xE0\x29\x3D\x76\x3A\x25\x92\x37\xFB" +
	"\xB2\x21\x3E\xC6\x9C\x57\xD6\xB9\xB0\x55\xD7\x25\x70\x7F\x5A\x9A" +
	"\x3B\xD9\xAD\x71\xB5\x9B\xCB\x3D\x8F\x22\xB7\x8E\x29\x37\xAA\x39" +
	"\x67\x97\xC9\x2B\xC5\x9B\x85\x69\x0A\xD5\x18\xF5\x74\x65\xCB\xC7" +
	"\x9F\x56\x43\xFD\x2A\xEC\x33\xC3\x71\xCC\x6E\x0F\xA8\xEE\x3F\x0A" +
	"\xDE\x35\x23\x2D\x99\xC5\x3A\x35\x29\xFC\x48\x43\x19\xDA\x5B\x1C" +
	"\x0E\xA6\x98\x56\xAC\x32\xF2\x71\x9C\x53\x0A\xD5\x99\xDC\xAE\x56" +
	"\x98\x45\x58\x65\xA8\xCA\xFB\x50\x52\x64\x0C\xA5\x58\x83\xC1\x14" +
	"\xA3\x1B\x08\x38\xCF\xF2\xA9\x19\x4E\xDF\x6C\xD3\x0A\xF1\x48\xD2" +
	"\xE3\xAD\xDF\xCA\x90\x1E\x30\x78\x39\xA2\x98\xCB\x8C\x83\xC1\x1D" +
	"\x8D\x15\x0E\x29\x82\x9B\x42\x8A\x70\x14\xD5\x15\x22\x8E\x2A\x8C" +
	"\xD8\xE5\x00\xF5\x34\xE0\x28\x51\x52\xED\x51\x18\x62\xE3\x3D\xC1" +
	"\xEC\x3D\x69\x92\x20\x15\x1C\xF7\x10\xDB\x2E\xE9\x9C\x2F\xB7\x73" +
	"\x54\x2F\x75\x65\x8B\x29\x6E\x41\xF5\x72\x38\xFC\x3D\x6B\x16\x20" +
	"\xD3\xCC\x59\xC9\x76\x3E\xA6\xB1\x9D\x65\x15\xA1\xD5\x47\x0D\x29" +
	"\xBD\x74\x36\x2E\x75\x29\x99\x33\x6C\x16\x35\xEB\x93\xCB\x63\xFA" +
	"\x55\x13\x28\x92\x5D\xC5\x9E\x59\x08\xEA\x6A\x45\x84\x92\x77\x60" +
	"\x27\x65\x14\x47\x6D\x14\x67\x70\xEB\xF5\xAE\x19\xD6\xE6\xDD\x9E" +
	"\xBD\x2C\x32\xA7\x6E\x54\x2F\x3B\x80\x18\x23\xBF\xB5\x0D\x0A\x3B" +
	"\x6E\x65\xC9\xA9\x82\x85\x20\x05\xEB\xE9\xDA\xAB\xB5\xBC\xAD\x26" +
	"\x4C\x80\xAF\xA5\x64\x9F\x99\xD4\xD7\x95\xC9\x12\x34\x51\xC0\x14" +
	"\xA5\x5B\x8D\x98\x03\xBD\x3C\x45\x94\xDA\x78\x1E\xD5\x20\x5C\x00" +
	"\x07\x6A\x97\x22\xAC\x56\x9A\x24\x7C\x79\x82\x95\x21\x44\x18\x5C" +
	"\xD2\xC8\x91\xAE\x0C\xCF\xC6\x72\x01\xA6\x1B\xB8\x54\x7C\xA4\x9F" +
	"\xA0\xA6\xAE\xD6\x84\xFB\xA9\xDD\x90\x35\xC2\xA4\x84\x08\xD8\x91" +
	"\xC6\x4D\x4F\x16\x48\xDC\x50\x2E\x7D\xE9\xA6\xEE\x16\xEB\x19\x3F" +
	"\x85\x2C\x6D\x13\x91\x88\xDB\xE8\x7A\x55\xBD\xB6\x26\x3A\xBD\xEE" +
	"\x48\x36\xB8\x2D\x90\x54\x75\x18\xA4\x8F\xCA\x6F\xF5\x64\x7E\x15" +
	"\x38\x5E\x38\x5C\x0A\x41\x12\x8E\x42\x01\x9A\xCB\x98\xD2\xC3\x02" +
	"\x00\x72\x07\x34\x63\x9A\x7A\xAB\x8C\xEE\xDB\xED\x8A\x44\x47\x00" +
	"\xEE\x21\x8F\x6E\x28\xB8\x0C\x2A\x39\xE2\xAB\x0B\x94\xDC\x55\x81" +
	"\x5A\x7A\xCB\x20\x94\xA4\xA3\xF2\xA9\x26\xB6\x49\x86\x7A\x13\xD1" +
	"\x85\x5A\xB2\xF8\x88\x77\x6A\xF1\x23\x11\xC6\xC7\x72\x75\xF6\xA7" +
	"\x04\x0A\xF9\x56\x2A\x47\xA1\xE6\x9B\x1C\x0D\x16\x01\x4C\xFF\x00" +
	"\xB4\xA6\xA4\x28\x1D\xB9\x18\x23\xA1\xA7\x7B\x3D\x18\xB9\x53\x5A" +
	"\xA2\xDD\xBE\xA6\xF1\x9D\x97\x0B\xE7\x28\xFE\x25\xFB\xC3\xEA\x2B" +
	"\x42\x19\xE1\xB8\x19\x89\xC3\x7B\x74\x3F\x95\x60\xEC\x8D\xDC\xB2" +
	"\x9F\x9B\xA1\xA2\x43\xB3\x0C\x41\xE3\xF8\x97\xA8\xAE\x98\x62\x64" +
	"\xB4\x7A\x9C\x15\x70\x10\x95\xDC\x74\x3A\x06\x5F\x6A\x61\xEF\x8E" +
	"\xF5\x42\xCF\x53\x0C\xA1\x27\x3B\x80\xE0\x48\x3F\xAD\x69\xE0\x15" +
	"\x0C\x39\x07\x90\x7D\x6B\xB6\x33\x52\x57\x47\x93\x52\x94\xA9\xBB" +
	"\x48\x81\xC7\x6E\x71\x4C\x90\xE4\x0E\x3B\x60\xD4\xCC\x2A\x39\x14" +
	"\xAE\x41\x1C\xD5\x09\x32\x03\x45\x0D\xC1\xA2\x81\x8F\x41\x9A\x91" +
	"\x45\x46\xB4\xE7\x95\x21\x8C\xBC\x8C\x15\x45\x2B\x93\x6E\x83\xE5" +
	"\x9A\x3B\x78\x8C\x92\xB0\x55\x15\x8F\x35\xEC\x97\xE1\x95\x4F\x97" +
	"\x08\xED\xDD\xBE\xB5\x56\xFA\xE9\xEE\xE4\x23\x04\xA9\x3F\x22\xFA" +
	"\x0A\xB1\x6F\x06\xD5\xF9\x94\x02\x7B\x57\x1D\x6A\xDA\x59\x1E\x9E" +
	"\x13\x0A\xAF\x79\xA1\xAB\x00\x71\xBA\x6C\x36\x0F\x07\xA7\x15\x3A" +
	"\xEC\x54\x05\x00\x0B\x9C\x70\x2A\x29\xA2\x9A\x59\x02\x81\x88\xC1" +
	"\xF5\xEB\x56\xD2\x00\x18\x12\x73\x8E\x83\xB0\xAE\x29\x4B\xBB\x3D" +
	"\x58\xC5\x2D\x90\x9B\x69\x22\x81\x23\x5F\x94\x71\xD7\x26\xAB\xDD" +
	"\x5C\x27\x98\xAB\x13\x33\x30\xEC\x3A\x1A\xB2\x91\xBC\x68\x8F\x2E" +
	"\xF7\x66\xE0\x46\xA3\x8A\x96\x9A\x5E\xA5\x26\x9B\x1E\x31\x8C\x83" +
	"\x90\x2A\x37\x95\x22\x20\x3E\xE0\x4F\x4E\x2A\xD9\xDB\x1A\x65\xF8" +
	"\xF6\x15\x59\xE0\x92\xE6\x60\x08\x21\x7D\x3D\x07\xBF\xBD\x44\x5A" +
	"\x7B\xEC\x50\xF8\xB1\x22\x06\x5E\x86\x9F\xB7\x04\x0D\xA4\xE7\xBF" +
	"\xA5\x3A\x45\x86\x38\xC4\x5E\x66\xCC\x74\x03\x93\x56\x02\xE5\x41" +
	"\xC1\xFC\xAA\x1C\xBA\x81\x81\x79\x65\x71\xE6\xB3\x80\x64\x04\xF1" +
	"\x8E\xD5\x4F\x63\x2F\x55\x23\xEA\x2B\xAC\x55\x04\x64\x54\x72\x18" +
	"\xC3\x6D\x65\x2C\x46\x01\xE3\xA5\x6B\x1C\x4B\xDA\xC6\x12\xA2\x9B" +
	"\xB9\xCF\x5B\xAA\x16\x05\xA3\x24\x7D\x6B\x65\x61\x46\x8C\x0D\x83" +
	"\x69\xED\x53\x2D\xBC\x32\x31\x65\x52\x1B\xE9\x8C\x54\x82\x02\x3A" +
	"\xBB\x11\xEF\x51\x52\xAA\x91\xBC\x52\x8C\x6C\x57\x58\x82\x8C\x01" +
	"\xC5\x35\x82\x42\x84\xB3\x60\x75\xE4\xE6\xAC\x88\x17\x9E\x59\xB3" +
	"\xC1\xE6\x9F\xE5\x00\x31\x81\x81\x59\xF3\x81\x90\xCB\x24\xD7\x01" +
	"\xAD\xF2\xAB\xFC\x4D\xD8\xD5\xB6\x42\x10\x90\x37\x10\x3A\x0E\xF5" +
	"\x66\x53\x1C\x11\x97\x90\x85\x51\x54\xC6\xA5\x6A\x5B\x19\x7F\xFB" +
	"\xE6\xB4\xBC\xA7\xB2\x15\xD2\x2A\xFD\xA8\x92\x41\xB7\x3C\x7E\x95" +
	"\x61\x1B\xCD\x5E\x15\x90\xFB\x8A\x9E\x47\x12\x42\x5A\xDC\x24\xAD" +
	"\xE9\x55\xE1\x97\xCC\xC8\x92\x37\x5F\x5D\xA7\x3F\xA5\x55\xEE\xAE" +
	"\x90\x24\x56\xD9\x77\x1C\xAA\xAC\xE0\x82\x7A\x9E\x95\x68\x02\x40" +
	"\xDC\x00\x35\x70\x20\x60\x08\xE4\x53\x0C\x2A\x0E\x42\xF3\x52\xEA" +
	"\x5C\x12\xB1\x49\xA1\x05\xB7\x2F\x0D\xEB\x49\xC3\x74\x20\xD5\x9F" +
	"\x20\x21\x62\x33\xF3\x75\x15\x4D\xAD\x0C\x52\x17\x4C\xED\x3D\x40" +
	"\xEA\x2B\x48\xC9\x3E\xA0\x35\xE3\x44\x05\x80\xDB\xFD\x69\xD6\x7A" +
	"\x83\xC2\x0A\xAF\xCE\x83\xAA\x31\xE9\xF4\xA5\x7C\x98\x8E\x06\xF2" +
	"\x07\x4E\x99\xA8\x05\xBC\x72\x2E\x57\x8C\xFE\x95\xB4\x26\xE3\xA9" +
	"\x85\x5A\x4A\x6B\x95\xAD\x0E\x86\x39\x60\xB8\x0A\xF1\x33\x18\xF1" +
	"\xF3\x64\x72\x0F\x7A\x59\xD7\x91\x95\x0A\x71\xDB\xA7\xD6\xB9\xFB" +
	"\x6B\xA9\x6C\xA7\x04\x2E\x71\xC3\x2F\xF7\x87\xB7\xBD\x76\x5A\x65" +
	"\xBE\x9D\xAA\x5B\xF9\xB6\xB3\xBB\x11\xF7\xE3\x6E\x19\x7E\xB5\xDD" +
	"\x1A\xC9\xAB\xB3\xC5\xAB\x86\x94\x25\x65\xB1\x84\xE2\x8A\xE8\xA5" +
	"\xD0\x90\x9F\x96\x4C\x0F\x7A\x2A\xBD\xB4\x4C\xFD\x94\x8E\x6D\xE6" +
	"\x48\x50\xB3\xB0\x02\xB1\xEE\x2E\x0D\xCC\xC5\xDD\x4F\x96\xA3\xE5" +
	"\x1E\xD5\x04\xF3\x35\xCC\xA5\xDF\x21\x07\xDD\x1E\x95\x3A\x3E\xD8" +
	"\xB7\x6D\xCF\x60\x3B\xD6\x15\x6A\x37\xA2\x3B\xF0\xB8\x74\xBD\xE9" +
	"\x6E\x4F\x6E\xB9\x45\x62\xBB\x4D\x4C\x24\x4F\x35\x63\x1F\x33\x1E" +
	"\xC3\xB5\x35\x44\x86\x26\x2B\xF7\xC8\xE0\x1E\xD4\xDD\x36\x09\x91" +
	"\x8E\xE8\x4E\xE2\x79\x66\x35\xC4\xED\x66\xD9\xEA\x27\x66\x91\x64" +
	"\x96\x57\xC0\x0A\xDF\xEC\xE7\x9A\xAF\xAA\xB4\xA8\x12\x35\x0C\x11" +
	"\x87\x2C\x3F\x95\x4F\x24\xB0\xDA\xCE\x63\x8D\x43\xCC\x4E\x4E\x7A" +
	"\x2D\x5E\x6C\xCA\x88\x91\x9C\xCA\x40\x24\x8E\x8B\xEE\x6B\x2E\x6E" +
	"\x56\xA5\x61\xC9\x73\x45\xD8\xC1\xD3\xE1\x97\xED\x0A\xF1\xC2\x64" +
	"\xC7\xAF\x02\xBA\x34\x47\x20\x17\x01\x3D\x40\xE6\xA6\x3B\x63\x51" +
	"\x9C\x00\x4E\x38\x1D\xE9\xE0\x56\x35\x6B\x39\xBB\xD8\x50\x8F\x2A" +
	"\xB1\x58\x5B\x2E\xE0\xC7\x93\xDF\xDE\xA4\xF2\x80\x04\x03\x8C\x8C" +
	"\x71\xDA\xA4\x74\x25\x48\x04\xAF\x1F\x78\x76\xAA\xE6\xF6\xDA\x32" +
	"\xB1\x99\x83\xBF\x4F\x97\x9A\xCD\x73\x4B\x62\xC8\x2D\x74\xE3\x6D" +
	"\x2B\x38\x65\x90\xB1\xC9\x67\x1C\xD5\x95\x84\xA9\x6C\x39\x19\xF5" +
	"\xE6\x89\x9E\x58\xAE\x17\x2C\x9E\x53\xF0\xA0\x82\x0E\x7E\xB5\x65" +
	"\xA3\x0D\xF7\x86\x68\x94\xE4\xF5\x7D\x41\x15\x25\x86\x5D\x8A\x22" +
	"\x97\x04\x1C\xB1\x6E\xAD\x4F\x2B\x20\x41\x85\x56\x6C\xF3\x93\xDA" +
	"\xA5\x9D\xBC\x98\x59\xC2\x33\xED\x1F\x75\x7A\x9A\x6C\x12\x34\x81" +
	"\x4B\xA8\x1B\xBA\x01\x9A\x57\x6D\x5C\x05\x00\x91\xF3\x0C\x1A\x69" +
	"\x0A\xC0\xAE\x41\xEC\x40\x35\x63\x6D\x35\x61\x54\xCE\xD1\x8C\x9C" +
	"\x9A\x9B\x85\x8A\xF1\xDB\x88\xC0\x0A\xC4\x01\xDB\xD6\xA4\xDB\x49" +
	"\x75\x71\x0D\xAA\x6E\x99\xB6\x83\xD0\x77\x35\x9E\x35\xC8\x59\x80" +
	"\x58\x64\x3F\x95\x5A\x84\xE7\xAA\x42\x6D\x2D\xC5\xD5\xAC\xE4\xB9" +
	"\xB6\x1E\x57\x2C\x87\x3B\x7D\x6B\x36\xCE\xD0\xBA\x00\xBF\x2C\x99" +
	"\xE8\xC3\xAF\xE3\x5A\xE9\xAA\xDB\x33\x00\xC2\x44\xCF\x72\xB5\x79" +
	"\x44\x73\x2E\x57\x0C\x3D\x45\x6A\xAA\xCE\x9C\x79\x5A\x08\xA5\x7E" +
	"\x64\x55\x4B\x68\xD0\x87\xD8\x11\x87\xA5\x35\x22\x59\x40\x91\x94" +
	"\x06\xEC\x47\x15\x74\x42\xB8\xC6\x32\x3D\xF9\xA4\x92\x32\xD1\x95" +
	"\x56\xD8\xC4\x70\x47\x6A\xC7\x9C\xA2\x01\x16\x31\xC9\x35\x04\x92" +
	"\xA2\x4E\xB0\x9C\xEE\x6E\x9C\x71\x56\xE1\x8E\x55\x4C\x4A\x43\x11" +
	"\xFC\x43\xBD\x0F\x1A\xB6\x37\x28\x38\xE4\x7B\x50\xA5\x67\xA8\x8A" +
	"\x92\x43\xB9\x08\x07\x19\xEE\x3B\x54\x4F\x09\x31\x15\xDC\x77\x63" +
	"\x1B\xAA\xCD\xD1\x68\xAD\xA4\x78\xC0\x2E\xA3\x20\x56\x12\xEA\xD7" +
	"\x01\x86\xF0\x8C\x3B\x8C\x62\xB7\xA7\x19\x4D\x5D\x11\x29\xA8\xEE" +
	"\x11\xB4\xF0\xCA\x62\x97\x2C\x73\xF8\xE3\xD4\x55\x9D\xA0\x8D\xC3" +
	"\xA1\xA7\xC5\x77\x6D\x79\x85\x90\x79\x6F\xDB\x3F\xD0\xD4\x8F\x13" +
	"\x46\x8E\x41\xDE\xC7\x90\x0F\x1F\x85\x6B\x29\x6B\xAA\xB3\x1C\x5A" +
	"\xB6\x86\x74\xD9\xDF\xB7\x68\xE9\x95\x34\xCB\x5B\xD9\x2D\x66\x13" +
	"\xDB\x16\x8E\x65\xE0\x81\xDC\x55\x99\x76\xAA\x2B\x48\x0F\x5E\x06" +
	"\x3A\x1A\xAB\x22\x2B\x92\xC3\x00\x9F\x4E\xF5\xBC\x24\x45\x48\x73" +
	"\x1D\x22\xEB\x33\x4F\x02\xB4\x73\xB1\x6C\x7C\xDC\xF4\x3E\x94\x57" +
	"\x29\x14\xB2\x40\xC2\x44\x38\x61\xC1\x1F\xDE\xFA\xD1\x5D\xD1\x9C" +
	"\x6D\xAA\x3C\x6A\xB4\x26\xA5\xEE\xBD\x08\xE3\x66\x07\x00\x06\xED" +
	"\x8A\xB9\x0A\x9D\xDB\x8F\x04\xF6\xEC\x2A\x9C\x0B\xFC\x59\xC0\x15" +
	"\x7A\xDF\x9C\x1D\xD9\x19\xE2\xB9\x2A\x1E\xAD\x15\x7D\xCB\x91\xE5" +
	"\x48\x38\x52\xBD\xCE\x71\xCD\x68\x85\x3B\x72\x9B\x73\xEF\xD2\xB1" +
	"\xC6\x64\xBD\x11\xCA\xAF\xE5\x2F\x2A\x14\x75\x35\xB2\x50\xBF\xCA" +
	"\x48\x58\xFD\x07\x7A\xE2\xAA\xAD\x63\xA2\xF7\x2A\xC7\xA6\x5B\xCF" +
	"\x29\xB8\x7D\xC4\xB1\xC9\x5C\xF1\x9A\xD3\x8D\x15\x10\x22\x28\x55" +
	"\x1D\x00\xA4\x51\xD8\x0A\xAF\xA9\xDE\x1B\x1B\x5D\xEA\x01\x76\x38" +
	"\x50\x7A\x7D\x6B\x16\xE5\x51\xA8\x8B\x48\xAB\x96\x9C\x46\xA5\x5D" +
	"\xC8\x04\x74\xC9\xA3\xCD\x89\x58\x29\x6C\x64\xE0\x64\x71\x58\x7A" +
	"\x4C\x53\x6A\x17\x3E\x64\xD2\x33\xC6\x0E\x49\xCF\x53\x5D\x1C\xD0" +
	"\x83\x6E\xCB\xB4\x91\x8C\x00\x3F\xA5\x15\x20\xA0\xD4\x5B\x2A\x37" +
	"\x71\xB9\x95\xAB\xC1\x27\xD9\x72\xB7\x22\x34\xCE\x5B\x7B\x75\xF6" +
	"\x15\x99\xA4\x14\x4B\x9D\xA8\x03\x93\xD5\xBA\x62\x99\xAF\x2C\xD1" +
	"\xDF\x84\x90\x9D\xAA\xA0\x2E\x7A\x56\xAE\x85\x6C\x89\x68\xCD\x13" +
	"\xC5\x2C\xE4\xE5\x41\x3D\x2B\xA9\x2E\x4A\x3A\xEB\x73\x38\xBE\x6A" +
	"\x8F\xC8\xB7\x88\x9A\x50\xD7\x0E\x41\x5E\x55\x5F\x81\xF5\x1E\xB5" +
	"\x0F\xDA\xA6\x7B\xC9\x63\x89\xA3\x65\x0A\x0A\x0E\xC7\xDF\x34\x97" +
	"\x2A\xD7\xF7\x50\xDB\xC6\xBB\x9A\x17\xDD\x33\xE3\xE5\x1E\xA1\x4D" +
	"\x6B\x8B\x58\x82\x85\xD8\xBB\x47\x40\x6B\x91\xB5\x14\xB9\xB7\x34" +
	"\x4A\xFB\x14\xE7\xB9\x8E\xD6\xDC\x4B\x3B\x01\x91\xFC\x3D\xCF\xB5" +
	"\x66\xA6\xBA\x65\x97\x64\x36\x8C\xE7\xEB\x57\x35\xCD\x32\x5B\xC8" +
	"\x11\xAD\xF9\x68\xBA\x27\xA8\xAC\xFD\x08\xC9\x65\x2B\xC5\x3D\xB4" +
	"\xC1\x9C\xF2\x76\x74\xFC\x7D\x2B\x4A\x70\x83\xA6\xE5\xBB\xEC\x4C" +
	"\x9C\x94\xED\xD0\xDA\x46\x97\x60\x69\x22\xC1\xC6\x48\x53\x9A\x95" +
	"\x40\x60\x08\xE4\x1A\x79\x71\xBB\x6A\x2B\x31\xEC\x71\xC6\x7E\xB4" +
	"\xE8\xA2\x65\x8C\x07\x20\xB7\x72\x07\x15\xCA\xFB\x9B\x58\xE3\x75" +
	"\xB5\x99\x2F\xDF\xCD\x25\x81\xFB\xA7\xB6\x2A\xDE\x9F\xA5\x48\xD2" +
	"\x2E\x4A\x85\xDA\x19\x9B\xBE\x0F\x61\x5D\x05\xE6\x99\x05\xE2\x6D" +
	"\x91\x71\xDF\x23\xD6\x92\x1B\x06\xB7\x88\x22\xCD\x23\x84\x18\x50" +
	"\x7A\xE3\xD2\xBA\xFE\xB2\xB9\x14\x56\x8C\xCA\x34\xED\x26\xD9\x1C" +
	"\x76\x76\xF1\x01\x18\x8F\xDC\x66\xA5\x54\x8C\x03\x12\x80\xB8\xEC" +
	"\x38\xA7\x89\x36\xEE\xCC\x52\x0C\x74\xC8\xEB\x55\x2F\x2F\x9D\x21" +
	"\x26\xDA\x26\xF3\x31\x91\xBD\x09\x07\xD8\x1F\x5A\xE7\x51\x9C\xD9" +
	"\x6D\xA4\x32\x7B\xCF\xB1\x5D\xC7\x0C\xE7\x74\x72\xFD\xD7\xEE\x0F" +
	"\xBD\x5C\x90\x30\x8D\x8A\x80\x5B\x07\x1F\x5A\xC4\x6B\x3D\x4B\x55" +
	"\x96\x33\x71\x6E\x2D\xE3\x07\x2C\xD8\xC1\xAE\x93\xCA\xC2\x80\x3B" +
	"\x0C\x53\xAB\x18\xC2\xDD\xFA\x93\x06\xE5\x73\x8C\x8E\xE2\xF6\xEA" +
	"\xE4\x85\x91\xF7\xE7\xA0\x38\xC5\x6A\xC4\x2F\xE3\x19\x97\x6C\x83" +
	"\x1D\x3B\xFE\x75\xA1\x2E\x9A\x9F\x69\x17\x30\x62\x39\x47\x5E\x38" +
	"\x6F\xAD\x48\xC8\xDB\x79\x52\x0F\xA5\x6B\x3A\xB1\x95\xAC\x82\x30" +
	"\xB6\xE5\x21\x20\x61\x91\xF8\x83\xDA\xB2\x2F\xF4\xCD\xCC\xD2\xDB" +
	"\x8E\x4F\x25\x3F\xC2\xB6\x8C\x2D\xBC\xB1\x5D\xBC\x63\xEB\x50\x5C" +
	"\xC0\xF2\x47\x84\x76\x46\x07\x20\x83\xD6\xAA\x9C\xB9\x5D\xD3\x26" +
	"\x50\xBA\xD4\xE7\x92\x01\x22\x02\x41\x56\x1C\x62\xAD\xC1\x34\x90" +
	"\xE1\x59\x8B\xC7\xE8\x7A\x8A\xD1\x7B\x65\x95\x72\xC0\xAB\xF4\xDD" +
	"\xEB\x59\xB7\x30\xCB\x03\x1C\xA3\x32\x67\xEF\x2F\x39\xAE\x8E\x75" +
	"\x3D\x19\xB4\x7D\x9A\x8E\xAB\x52\xC5\xC4\x62\x45\xC6\x3E\x52\x32" +
	"\x0D\x66\xED\x2A\xE4\x11\xF5\xAB\x76\x72\x34\xCB\x2C\x5F\x38\x00" +
	"\x64\x12\x3A\x53\x26\x07\x77\x7E\x3F\x5A\x23\x78\xBE\x56\x45\xD3" +
	"\xD8\xA6\xEC\x08\x6E\xE5\x4F\x7A\x29\xAF\x8F\x34\xF7\xCF\xF3\xA2" +
	"\xBA\x11\xCA\xDB\x6F\x72\x22\xF8\x41\xC5\x5A\x86\x46\xC2\x6D\x4C" +
	"\x9E\xF8\xED\x54\xCE\x36\x81\xEA\x73\x56\xE1\x2E\x91\x0D\x83\x2C" +
	"\x7F\x4A\x26\xB4\x15\x36\xEE\x6D\xD8\xEF\xF2\xB3\x21\xE4\x9E\x9E" +
	"\x95\x3C\xF3\x2C\x0A\xA5\x81\x3B\x8E\x00\x03\x35\x5E\xD5\xB6\x44" +
	"\x37\x1E\x7B\x9A\xB4\xCC\xDE\x5B\x18\xC0\x67\xC7\xCA\x0F\x4C\xD7" +
	"\x9D\x25\xEF\x1D\x05\x98\x41\x60\x09\x18\x34\xEB\xCD\x3A\x3B\xFB" +
	"\x6F\x2A\x42\x47\x75\x61\xD8\xD5\x7D\x3D\xA4\xF2\xFF\x00\x7E\xEA" +
	"\xD2\xF7\x03\xB5\x6A\x23\x0C\x56\x52\xBC\x25\x74\x68\x92\x92\xD4" +
	"\xC6\xD3\xF4\x2B\xAB\x29\x89\x8E\xF8\x2C\x67\xA8\x0B\x92\x7F\x3A" +
	"\xDC\x8E\x00\x9B\x49\x67\x76\x1D\xD8\xF7\xF5\xA6\xBD\xC4\x71\x32" +
	"\x23\xB0\x0E\xE7\x0A\xBD\xCD\x49\x0C\xA2\x54\xDC\x06\x06\x48\xA8" +
	"\xA9\x3A\x93\xD6\x44\xA4\xA3\xA2\x07\x82\x39\x46\x24\x8D\x1C\x7F" +
	"\xB4\x33\x55\x53\x4C\xB0\xF3\x84\xB1\xC2\x81\xBA\x7C\xA7\x83\x55" +
	"\xB5\xA1\x7D\x24\xD6\xEB\x6E\x8C\xF6\xB9\xFD\xEA\xA1\xC1\x3C\xF7" +
	"\xF6\xAD\x09\x21\x0F\x1A\xC0\x8B\xE5\xC6\x31\x96\x1C\x10\x3D\xA8" +
	"\x49\xC6\x29\xF3\x6E\x2B\xDF\xA0\x5C\x49\x15\x8D\x9B\xCA\xE3\x6C" +
	"\x51\x8C\x90\xA2\xB9\x61\xAD\xEA\x5A\x8D\xE7\x95\x66\x04\x2A\x79" +
	"\xC2\x8C\x90\x3D\x49\xAE\xB6\xEE\x08\xEE\xED\xA4\x82\x5C\xEC\x90" +
	"\x60\xE3\xB5\x60\x69\xBA\x14\xFA\x76\xA0\x64\x26\x39\xE0\x23\x1C" +
	"\x64\x35\x6B\x87\x95\x34\x9B\x96\xFE\x64\x49\xC9\xB4\x96\xC6\xC5" +
	"\x9A\x5C\x2D\xAA\x7D\xA6\x55\x96\x66\xC9\x04\x70\x0F\xA5\x4F\x6E" +
	"\xE6\x58\x86\x58\x17\x1F\x7C\x0E\xC6\x96\x47\x90\x90\x22\x40\xB9" +
	"\x1F\x79\x8F\xDD\xFC\x2A\xA3\x5F\x2D\xAD\xE4\x7A\x7D\xBD\xB3\xCD" +
	"\x21\x1B\x89\x2D\x81\x8E\xE7\x3D\xEB\x1E\x57\x3B\xDB\x72\xF9\x9A" +
	"\x2F\x6D\xA6\xB4\x72\xF9\x8B\xB0\x2F\x97\xFC\x44\xF5\xAA\x1A\xF6" +
	"\xA5\x2D\x86\x9B\xE7\x5B\x2E\x5D\x9B\x66\xE6\x1C\x2F\xBD\x73\x9A" +
	"\x75\xC6\xAD\xA8\xCA\xEF\xF6\xCB\x85\x8D\x01\x66\x75\xE7\xF0\x02" +
	"\xAE\x96\x1E\x53\x8F\x3B\x76\x44\x3A\x9A\xF2\x9D\xA1\x43\x83\x8C" +
	"\x67\xB6\x68\x54\x6D\xA3\x76\x0B\x77\xC5\x54\xD3\x7E\xD0\x51\x64" +
	"\x37\x2F\x34\x6C\xA0\xED\x95\x46\x41\xFA\x8F\xE5\x57\x2D\xEE\x22" +
	"\xB9\x46\x68\xF7\x0D\xAC\x55\x95\x86\x08\x22\xB1\x9C\x1C\x6E\xB7" +
	"\x29\xC9\x86\xDA\x50\xA7\xB5\x49\x91\x55\xEE\xC4\x6E\xAB\x13\xB3" +
	"\xAE\xE3\x90\x54\xFA\x54\x45\x5D\xD8\x14\x98\xF2\xB5\x1C\x88\x80" +
	"\xE5\xD8\x2E\xE1\x8E\x5B\x15\x30\x00\x0E\x1B\x71\x51\x8F\xA9\xAF" +
	"\x39\xBB\x6B\xDB\xFD\x4D\xE3\x93\x73\x4C\x5C\x8D\xA4\xF0\x2B\xA3" +
	"\x0F\x43\xDA\xB7\xAD\xAC\x4C\xAA\x38\x9E\x85\xB3\x70\xC8\x20\x8F" +
	"\x50\x73\x4C\x68\xEB\x2B\x46\xD1\x9A\xCD\x16\x57\x96\x41\x2F\x75" +
	"\x57\xE2\xB5\x56\x70\xD3\xCB\x09\xDA\x59\x00\x3F\x29\xEC\x7D\x7D" +
	"\x0D\x4C\xE0\xA2\xED\x17\x73\x4E\x66\xB7\x29\xDF\x59\xC9\x3A\x2F" +
	"\x95\x33\x44\xCA\x73\xC7\x7A\x86\xE0\xAD\xAD\xB7\x99\x70\xC0\x28" +
	"\xC0\x27\xDE\xB2\xA4\xF1\x05\xFC\x3A\x89\x82\xE2\x28\x63\x5D\xF8" +
	"\xC3\x0C\x60\x7D\x6B\x76\x4F\x2E\xEA\xDD\x91\xD7\x28\xC3\x90\x6B" +
	"\x69\x42\x74\xD2\xE6\xD8\x23\x2E\x6B\xD8\xCE\x8E\x68\x6E\x10\xBC" +
	"\x32\x2B\x8F\x6E\xD5\x52\xE2\xFA\x0B\x69\x76\x49\x20\x0D\xE9\xE9" +
	"\x52\x41\x61\xF6\x09\x9F\xC9\x2A\xD1\x37\x76\xFB\xC3\xFC\x6B\x9F" +
	"\xD6\xE3\x29\xA8\x33\xF5\x12\x0C\x8A\xE9\xA7\x4E\x33\x95\x93\xD0" +
	"\xCA\xAC\xE5\x08\x5E\xC6\xBB\xCE\xEC\xEB\x24\x05\x25\x84\xF0\xC0" +
	"\x75\x1E\xF5\x97\x78\xAF\x15\xC3\x32\xB9\x31\xB9\xE7\xDA\xA8\xDB" +
	"\xCF\x25\xBC\x9B\xA3\x38\xF5\x1D\x8D\x68\x49\x70\xB2\x44\x24\x44" +
	"\xC9\x27\x0C\x3D\x2B\x75\x4D\xD3\x7A\x6C\x65\x1A\x8A\xA2\xD7\x72" +
	"\xB3\xFD\xD3\xB4\xF3\x8E\x0F\xAD\x14\x4A\xC4\x20\x60\x3F\x03\x45" +
	"\x6D\x1B\x91\x52\xD7\xDC\xFF\xD9"
