package models

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/nfnt/resize"
	"golang.org/x/net/context"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

func Upload(id string, file multipart.File) (message string, err error) {
	log.Print(id)
	var pos = strings.Split(id, ".")
	log.Print(pos[len(pos)-1])

	switch pos[len(pos)-1] {
	case "png":
		return makePng(id, file)
	case "jpeg":
		return makeJpeg(id, file)
	case "jpg":
		return makeJpeg(id, file)
	default:
		log.Print(pos[len(pos)-1])
	}
	return makePng(id, file)
}

func makePng(id string, file multipart.File) (message string, err error) {
	img, _ := png.Decode(file)
	buf := new(bytes.Buffer)
	png.Encode(buf, img)

	var outRect image.Image
	outRect = ImageResize(img)

	conf, _, err := image.DecodeConfig(buf)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("Width=%d, Height=%d\n", conf.Width, conf.Height)

	out, err := os.Create("tmp.png")
	if err != nil {
		log.Print(err)
	}
	defer out.Close()
	err = jpeg.Encode(out, outRect, nil)
	if err != nil {
		log.Print(err)
	}

	f, err := os.Open("tmp.png")
	if err != nil {
		log.Print(err)
	}
	defer f.Close()

	return PutS3(id, f, "png")
}

func makeJpeg(id string, file multipart.File) (message string, err error) {
	img, _ := jpeg.Decode(file)
	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)

	conf, _, err := image.DecodeConfig(buf)
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("Width=%d, Height=%d\n", conf.Width, conf.Height)

	var outRect image.Image

	outRect = ImageResize(img)

	out, err := os.Create("tmp.jpeg")
	if err != nil {
		log.Print(err)
	}
	defer out.Close()
	err = jpeg.Encode(out, outRect, nil)
	if err != nil {
		log.Print(err)
	}

	f, err := os.Open("tmp.jpeg")
	if err != nil {
		log.Print(err)
	}
	defer f.Close()
	//defer file.Close()
	return PutS3(id, f, "jpeg")
}

func PutS3(id string, f multipart.File, pos string) (message string, err error) {
	ctx := context.Background()
	S3URL := beego.AppConfig.String("S3URL")
	accessKey := beego.AppConfig.String("accessKey")
	secretKey := beego.AppConfig.String("::secretKey")
	bucket := beego.AppConfig.String("::bucket")
	// Initial credentials loaded from SDK's default credential chain. Such as
	// the environment, shared credentials (~/.aws/credentials), or EC2 Instance
	// Role. These credentials will be used to to make the STS Assume Role API.
	sess := session.Must(session.NewSession())

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")

	// Create service client value configured for credentials
	// from assumed role.
	svc := s3.New(sess, &aws.Config{Credentials: creds, Region: aws.String(endpoints.ApNortheast1RegionID)})

	var contentType = "image/" + pos
	// Uploads the object to S3. The Context will interrupt the request if the
	// timeout expires.
	_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String("image/" + id),
		Body:        f,
		ContentType: &contentType,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == request.CanceledErrorCode {
			// If the SDK can determine the request or retry delay was canceled
			// by a context the CanceledErrorCode error code will be returned.
			fmt.Fprintf(os.Stderr, "upload canceled due to timeout, %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "failed to upload object, %v\n", err)
		}
		os.Exit(1)
	}

	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, "image/"+id)
	return S3URL + "/image/" + id, nil
}

func ImageResize(img image.Image) (out image.Image) {
	pngImg := resize.Resize(50, 0, img, resize.Lanczos3)
	// 書き出し用のイメージを準備
	outRect := image.NewRGBA(image.Rect(0, 0, 50, 50))
	//srcRect := image.Rectangle{image.Pt(0, 0), pngImg.Bounds().Size()}
	draw.Draw(outRect, image.Rect(0, 0, 50, 50), pngImg, image.Pt(0, 0), draw.Src)
	return outRect
}
