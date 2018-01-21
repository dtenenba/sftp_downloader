package main

import (
	"os"
	"testing"
	"time"

	"github.com/FredHutch/sftp_downloader/mocks"
	"github.com/golang/mock/gomock"
)

func getTheTime() time.Time {
	ret, _ := time.Parse(time.RFC822, time.RFC822) // 2006-01-02 15:04:00 +0000 MST
	return ret
}

func TestDoTheWork(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockSftper := mocks.NewMockSftper(mockCtrl)
	mockFiler := mocks.NewMockFiler(mockCtrl)
	mockWalker := mocks.NewMockWalker(mockCtrl)
	mockFileInfo := mocks.NewMockFileInfo(mockCtrl)

	mockSftper.EXPECT().Close().Return(nil).Times(1)
	mockSftper.EXPECT().Walk("/tmp/").Return(mockWalker).Times(1)
	mockSftper.EXPECT().Create("hello.txt").Return(mockFiler, nil).Times(1)
	mockSftper.EXPECT().Lstat("hello.txt").Return(mockFileInfo, nil).Times(1)
	mockWalker.EXPECT().Step().Return(true).Times(1)
	mockWalker.EXPECT().Step().Return(false).Times(1)
	mockWalker.EXPECT().Err().Return(nil).Times(1)
	mockWalker.EXPECT().Path().Return("/tmp/something").Times(1)
	mockFiler.EXPECT().Write([]byte("Hello world!")).Return(12, nil).Times(1)

	doTheWork(mockSftper)

}

func TestGetDateString(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	t.Run("garbage", func(t *testing.T) {
		os.Args = []string{"sftp_uploader", "config.json", "garbage"}
		_, err := getDateString()
		if err == nil {
			t.Error("Expected error with getDateString() and invalid date string")
		}
	})

	t.Run("yesterday", func(t *testing.T) {
		os.Args = []string{"sftp_uploader", "config.json"}
		currentTimeFunction = getTheTime
		ds, err := getDateString()
		if err != nil {
			t.Error("Did not expect error in getDateString()")
		}
		if ds != "01-01-2006" {
			t.Error("Expected 01-01-2006, got", ds)
		}
	})

}
