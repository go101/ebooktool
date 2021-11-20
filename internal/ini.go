package internal

import (
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strconv"

	"github.com/zieckey/goini"

	"go101.org/ebooktool/internal/nstd"
)

// For the ini package, to support
// * unkeyed values (under a section, as value list)
// * or allow duplicated keys to form lists
// * or, support object and list values
//      key1: {
//		k1: v1,
//		k2: v2,
//      }
//      key2: {
//		v1,
//		v2,
//      }

type Config struct {
	ini  *goini.INI
	path string
}

func LoadIniFile(iniFile string) (*Config, error) {
	data, err := os.ReadFile(iniFile)
	if err != nil {
		return nil, err
	}

	ini := goini.New()
	ini.SetParseSection(true)
	ini.SetSkipCommits(true)
	ini.SetTrimQuotes(true)

	err = ini.Parse([]byte(data), "\n", ":")
	if err != nil {
		return nil, err
	}

	err = mergeIncluded(ini, iniFile)
	if err != nil {
		return nil, err
	}

	//println(ini.Get("title"))
	//println(ini.Get("author"))
	//println(ini.Get("cover-image"))
	//println(ini.SectionGetInt("version-text-on-cover", "x"))

	return &Config{ini: ini, path: iniFile}, nil
}

const IniKeyInclude = "include"

func mergeIncluded(ini *goini.INI, includingFilepath string) error {
	includedFile, ok := ini.Get(IniKeyInclude)
	if !ok {
		return nil
	}

	inlcuded := goini.New()
	inlcuded.SetParseSection(true)
	inlcuded.SetSkipCommits(true)
	inlcuded.SetTrimQuotes(true)

	includedFilepath := filepath.Join(filepath.Dir(includingFilepath), includedFile)
	data, err := os.ReadFile(includedFilepath)
	if err != nil {
		return err
	}
	err = inlcuded.Parse([]byte(data), "\n", ":")
	if err != nil {
		return err
	}

	ini.Delete(goini.DefaultSection, IniKeyInclude)
	ini.Merge(inlcuded, false)
	return mergeIncluded(ini, includedFilepath)
}

func (c *Config) MakePath(relPath string) string {
	return filepath.Join(filepath.Dir(c.path), relPath)
}

func (c *Config) Path(keys ...string) (string, bool) {
	v, ok := c.String(keys...)
	if ok {
		v = c.MakePath(v)
	}
	return v, ok
}

func (c *Config) Color(keys ...string) (color.RGBA, bool) {
	v, ok := c.String(keys...)
	if ok {
		if nstd.String(v).HasPrefix("#") {
			v = v[1:]
		}
		n, err := strconv.ParseUint(v, 16, 64)
		if err != nil {
			return color.RGBA{}, true
		}
		return color.RGBA{
			A: 255, // byte((n >> 24) & 255),
			R: byte((n >> 16) & 255),
			G: byte((n >> 8) & 255),
			B: byte((n >> 0) & 255),
		}, true
	}
	return color.RGBA{}, false
}

func (c *Config) String(keys ...string) (string, bool) {
	switch len(keys) {
	default:
		nstd.Panicf("Config.Get only suports one and two arguments now, but got %d", len(keys))
	case 1:
		return c.ini.Get(keys[0])
	case 2:
		return c.ini.SectionGet(keys[0], keys[1])
	}
	panic("unreachable")
}

func confirmAnchor(x, y, z bool, nameList string) (int8, error) {
	var r int8 = 0
	var n = 0
	if x {
		r = -1
		n++
	}
	if y {
		r = 0
		n++
	}
	if z {
		r = 1
		n++
	}
	if n != 1 {
		return 0, fmt.Errorf("x anchor must be specified once in %s, but now %d are specified", nameList, n)
	}
	return r, nil
}

func (c *Config) CoordinateX(keys ...string) (x int32, anchor int8, err error) {
	switch len(keys) {
	default:
		nstd.Panicf("Config.Get only suports one and two arguments now, but got %d", len(keys))
	case 1:
		left, okLeft := c.ini.GetFloat(keys[0] + ".left")
		center, okCenter := c.ini.GetFloat(keys[0] + ".center")
		right, okRight := c.ini.GetFloat(keys[0] + ".right")
		if anchor, err = confirmAnchor(okLeft, okCenter, okRight, "left|center|right"); err != nil {
			return
		}
		x = int32(left + center + right)
		return
	case 2:
		left, okLeft := c.ini.SectionGetFloat(keys[0], keys[1]+".left")
		center, okCenter := c.ini.SectionGetFloat(keys[0], keys[1]+".center")
		right, okRight := c.ini.SectionGetFloat(keys[0], keys[1]+".right")
		if anchor, err = confirmAnchor(okLeft, okCenter, okRight, "left|center|right"); err != nil {
			return
		}
		x = int32(left + center + right)
		return
	}
	panic("unreachable")
}

func (c *Config) CoordinateY(keys ...string) (y int32, anchor int8, err error) {
	switch len(keys) {
	default:
		nstd.Panicf("Config.Get only suports one and two arguments now, but got %d", len(keys))
	case 1:
		top, okTop := c.ini.GetFloat(keys[0] + ".top")
		middle, okMiddle := c.ini.GetFloat(keys[0] + ".middle")
		bottom, okBottom := c.ini.GetFloat(keys[0] + ".bottom")
		if anchor, err = confirmAnchor(okTop, okMiddle, okBottom, "top|middle|bottom"); err != nil {
			return
		}
		y = int32(top + middle + bottom)
		return
	case 2:
		top, okTop := c.ini.SectionGetFloat(keys[0], keys[1]+".top")
		middle, okMiddle := c.ini.SectionGetFloat(keys[0], keys[1]+".middle")
		bottom, okBottom := c.ini.SectionGetFloat(keys[0], keys[1]+".bottom")
		if anchor, err = confirmAnchor(okTop, okMiddle, okBottom, "top|middle|bottom"); err != nil {
			return
		}
		y = int32(top + middle + bottom)
		return
	}
	panic("unreachable")
}

func (c *Config) Int32(keys ...string) (int32, bool) {
	v, ok := c.Number(keys...)
	if ok {
		return int32(v), true
	}
	return 0, false
}

func (c *Config) Number(keys ...string) (float64, bool) {
	switch len(keys) {
	default:
		nstd.Panicf("Config.Get only suports one and two arguments now, but got %d", len(keys))
	case 1:
		return c.ini.GetFloat(keys[0])
	case 2:
		return c.ini.SectionGetFloat(keys[0], keys[1])
	}
	panic("unreachable")
}

func (c *Config) Bool(keys ...string) (bool, bool) {
	switch len(keys) {
	default:
		nstd.Panicf("Config.Get only suports one and two arguments now, but got %d", len(keys))
	case 1:
		return c.ini.GetBool(keys[0])
	case 2:
		return c.ini.SectionGetBool(keys[0], keys[1])
	}
	panic("unreachable")
}
