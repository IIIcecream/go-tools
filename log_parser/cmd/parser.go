package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Parser struct {
	modules []string
	unzip   bool
	merge   bool
	srcZip  string
	destDir string
}

type ParserOptions struct {
	modules []string
	unzip   bool
	merge   bool
	srcZip  string
	destDir string
}
type ParserOption func(*ParserOptions)

func WithModules(modules ...string) ParserOption {
	return func(o *ParserOptions) {
		o.modules = modules
	}
}

func WithUnzip(unzip bool) ParserOption {
	return func(o *ParserOptions) {
		o.unzip = unzip
	}
}

func WithSrcZip(srcZip string) ParserOption {
	return func(o *ParserOptions) {
		o.srcZip = srcZip
	}
}

func WithDestDir(destDir string) ParserOption {
	return func(o *ParserOptions) {
		o.destDir = destDir
	}
}

func WithMerge(merge bool) ParserOption {
	return func(o *ParserOptions) {
		o.merge = merge
	}
}

func transCurDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	exeDir, err = filepath.Abs(exeDir)
	return exeDir, nil
}

func NewParser(opts ...ParserOption) (*Parser, error) {
	options := &ParserOptions{}

	for _, opt := range opts {
		opt(options)
	}

	fileName := filepath.Base(options.srcZip)                              // 获取文件名，如 "example.zip"
	nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName)) // 去掉扩展名，如 "example"
	dirPath := filepath.Dir(options.srcZip)                                // 获取目录路径，如 "D:\workspace\logs"
	options.destDir = filepath.Join(dirPath, nameWithoutExt)               // 拼接最终路径

	if options.destDir == "./" {
		if options.srcZip != "./" {
			options.destDir = dirPath
		} else {
			absDir, err := transCurDir()
			if err != nil {
				return nil, err
			}
			options.destDir = absDir
		}
	}

	return &Parser{
		modules: options.modules,
		unzip:   options.unzip,
		merge:   options.merge,
		srcZip:  options.srcZip,
		destDir: options.destDir,
	}, nil
}

type ParsedResult struct {
	filePath  string
	timestamp int64
}

func (parser *Parser) Parse() error {
	if parser.unzip {
		err := Unzip(parser.srcZip, parser.destDir)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// 遍历 parser.destDir 文件夹，提取要操作的文件
	files, err := parser.getFiles()
	if err != nil {
		return err
	}

	var errList []string
	var parsedResults []ParsedResult
	for _, inputFile := range files {
		outputResult, err := parser.convert(inputFile)
		if err != nil {
			errList = append(errList, fmt.Sprintf("%s: %v", inputFile, err))
			continue
		}
		parsedResults = append(parsedResults, outputResult)
	}

	if parser.merge {
		err := parser.Merge(parsedResults)
		if err != nil {
			errList = append(errList, fmt.Sprintf("merge failed: %v", err))
		}
	}

	var errs error
	if len(errList) > 0 {
		errs = fmt.Errorf("some files failed:\n%s", strings.Join(errList, "\n"))
	}
	return errs
}

// 获取解压后，提取出来的待处理的 log 文件列表
func (parser *Parser) getFiles() ([]string, error) {
	var ret []string
	err := filepath.WalkDir(parser.destDir,
		func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if strings.HasSuffix(d.Name(), ".rrzipped") {
				matched := false
				for _, module := range parser.modules {
					if strings.Contains(d.Name(), module) {
						matched = true
						break
					}
				}
				if !matched {
					return nil
				}

				ret = append(ret, path)
			}

			return nil
		})
	return ret, err
}

// 1. 第一行的 base + 每一行的 offset 进行时间戳还原
// 2. 将时区 +8
func (parser *Parser) convert(filePath string) (ParsedResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ParsedResult{}, err
	}
	defer file.Close()

	newFilePath := filePath + ".txt"
	outFile, err := os.Create(newFilePath)
	if err != nil {
		return ParsedResult{}, err
	}
	defer outFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	var baseTime time.Time
	var baseTimestamp int64

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		if lineNum == 1 {
			// 第一行：时间字符串 + timestamp
			// 示例：2025/12/25 03:11:37 28630693
			parts := strings.Fields(line)
			if len(parts) < 3 {
				return ParsedResult{}, fmt.Errorf("invalid first line format")
			}

			baseTime, err = time.Parse("2006/01/02 15:04:05", parts[0]+" "+parts[1])
			baseTime = baseTime.Add(8 * time.Hour) // +8 时区
			if err != nil {
				return ParsedResult{}, fmt.Errorf("parse base time failed: %v", err)
			}
			baseTimestamp, err = strconv.ParseInt(parts[2], 10, 64)
			if err != nil {
				return ParsedResult{}, fmt.Errorf("parse base timestamp failed: %v", err)
			}

			continue
		}

		// 后续行，按逗号或空格分割
		cols := strings.FieldsFunc(line, func(r rune) bool {
			return r == ',' || r == ' '
		})

		if len(cols) < 2 {
			// 不符合预期的行，直接写回
			writer.WriteString(line + "\n")
			continue
		}

		// 第二列时间戳
		ts, err := strconv.ParseInt(strings.TrimSpace(cols[1]), 10, 64)
		if err != nil {
			writer.WriteString(line + "\n")
			continue
		}

		// 计算时间差
		diff := ts - baseTimestamp
		newTime := baseTime.Add(time.Duration(diff) * time.Millisecond)

		cols[1] = newTime.Format("2006/01/02 15:04:05.000")
		newLine := strings.Join(cols[1:], " ")
		writer.WriteString(newLine + "\n")
	}

	if err := scanner.Err(); err != nil {
		return ParsedResult{}, err
	}

	return ParsedResult{
		filePath:  newFilePath,
		timestamp: baseTimestamp,
	}, nil
}

func appendFile(w *bufio.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}

func (parser *Parser) Merge(results []ParsedResult) error {
	// module -> []result
	grouped := make(map[string][]ParsedResult)

	// 1. 按 module 分组
	for _, result := range results {
		base := filepath.Base(result.filePath)
		for _, module := range parser.modules {
			if strings.Contains(base, module) {
				grouped[module] = append(grouped[module], result)
				break
			}
		}
	}

	var errList []string

	// 2. 逐 module 处理
	for module, moduleResult := range grouped {
		// 按时间排序
		sort.Slice(moduleResult, func(i, j int) bool {
			return moduleResult[i].timestamp < moduleResult[j].timestamp
		})

		// 3. 合并写文件
		outName := fmt.Sprintf("%s.log.merged", module)
		out, err := os.Create(parser.destDir + "\\" + outName)
		if err != nil {
			errList = append(errList,
				fmt.Sprintf("[%s] create output failed: %v", module, err))
			continue
		}
		defer out.Close()

		writer := bufio.NewWriter(out)
		defer writer.Flush()

		for _, item := range moduleResult {
			if err := appendFile(writer, item.filePath); err != nil {
				errList = append(errList,
					fmt.Sprintf("[%s] %s: %v", module, item.filePath, err))
			}
		}
	}

	// 4. 统一返回错误
	if len(errList) > 0 {
		return fmt.Errorf("merge finished with errors:\n%s",
			strings.Join(errList, "\n"))
	}

	return nil
}
