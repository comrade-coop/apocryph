import { outdent } from 'outdent'

export default {
  'Go': {
    language: 'go',
    code: outdent`
      package main

      import (
        "log"

        _ "github.com/joho/godotenv/autoload"
        "github.com/minio/minio-go/v7"
        "github.com/minio/minio-go/v7/pkg/credentials"
      )

      func main() {
        creds, _ := credentials.NewCustomTokenCredentials(
          os.Getenv("S3_BUCKET"),
          os.Getenv("CUSTOM_TOKEN"),
          "arn:minio:iam:::role/idmp-swieauth",
        )
        minioClient, err := minio.New(os.Getenv("S3_BUCKET"), &minio.Options{
          Secure: true,
          Creds:  creds,
        })
        if err != nil {
          log.Fatalln(err)
        }

        log.Printf("%#v\\n", minioClient) // Do your S3 magic! ✨
      }
    `
  },
  'JavaScript': {
    language: 'javascript',
    code: outdent`
      import 'dotenv/config'
      import * as Minio from 'minio'

      const minioClient = new Minio.Client({
        endPoint: process.env.S3_BUCKET,
        port: 9000,
        useSSL: true,
        credentialsProvider: new AssumeRoleWithCustomTokenProvider(
          process.env.S3_BUCKET,
          process.env.CUSTOM_TOKEN,
          "arn:minio:iam:::role/idmp-swieauth"
        )
      })

      minioClient // Do your S3 magic! ✨
    `
  },
}
export const envExample = (bucketLink: string, siweToken: string) =>
  outdent`
    # .env
    S3_BUCKET='${bucketLink}'
    CUSTOM_TOKEN='${siweToken}'
  `
