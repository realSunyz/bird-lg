package bird

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const responseCodeAccessRestricted = 16

func Query(socketPath, command string, timeout time.Duration) (string, error) {
	conn, err := net.DialTimeout("unix", socketPath, 5*time.Second)
	if err != nil {
		if isTimeoutError(err) {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("bird_connection_failed")
	}
	defer conn.Close()

	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	conn.SetDeadline(time.Now().Add(timeout))

	reader := bufio.NewReader(conn)

	if _, _, err = readResponse(reader); err != nil {
		if isTimeoutError(err) {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("bird_welcome_failed")
	}

	if _, err = conn.Write([]byte("restrict\n")); err != nil {
		if isTimeoutError(err) {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("bird_restrict_failed")
	}
	_, responseCode, err := readResponse(reader)
	if err != nil || responseCode != responseCodeAccessRestricted {
		if isTimeoutError(err) {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("bird_restrict_failed")
	}

	if _, err = conn.Write([]byte(command + "\n")); err != nil {
		if isTimeoutError(err) {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("bird_command_failed")
	}

	response, _, err := readResponse(reader)
	if err != nil {
		if isTimeoutError(err) {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("bird_response_failed")
	}
	return response, nil
}

func isTimeoutError(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func readResponse(reader *bufio.Reader) (string, int, error) {
	var result strings.Builder
	isEndLine := func(line string) bool {
		if len(line) < 5 {
			return false
		}
		for i := 0; i < 4; i++ {
			c := line[i]
			if c < '0' || c > '9' {
				return false
			}
		}
		return line[4] == ' '
	}

	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			if len(line) > 5 && line[4] == '-' {
				line = line[5:]
			} else if isEndLine(line) {
				responseCode, parseErr := strconv.Atoi(line[:4])
				if parseErr != nil {
					return result.String(), -1, parseErr
				}
				if len(line) > 5 {
					result.WriteString(line[5:])
				}
				return result.String(), responseCode, nil
			}

			result.WriteString(line)
		}
		if err != nil {
			return result.String(), -1, err
		}
	}
}
