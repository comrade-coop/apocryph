import { outdent } from 'outdent'

export default {
  'JavaScript': {
    language: 'javascript',
    code: (withToken: boolean) => outdent`
      import 'dotenv/config'
      import * as Minio from 'minio'

      const minioClient = new Minio.Client({
        endPoint: process.env.S3_BUCKET,
        port: 9000,
        useSSL: true,
        ${withToken ?
          `credentialsProvider: new AssumeRoleWithCustomTokenProvider(process.env.CUSTOM_TOKEN),` :
          `accessKey: process.env.ACCESS_KEY,\n  secretKey: process.env.SECRET_KEY,`
        }
      })

      minioClient // Do your S3 magic! ✨
    `
  },
  'JavaScript (AWS)': {
    language: 'javascript',
    code: (withToken: boolean) => outdent`
      import 'dotenv/config'
      import { S3Client } from "@aws-sdk/client-s3"

      const s3Client = new S3Client({
        endPoint: process.env.S3_BUCKET,
        region: "${'us-east-1' /* Minio Default */}",
        credentials: {${withToken ? `// The AWS SDK doesn't support custom tokens` : ''}
          accessKeyId: process.env.ACCESS_KEY,
          secretAccessKey: process.env.SECRET_KEY,
        },
      })

      s3Client // Do your S3 magic! ✨
    `
  },
  'Go': {
    language: 'go',
    code: (withToken: boolean) => outdent`
      package main

      import (
        "log"

        _ "github.com/joho/godotenv/autoload"
        "github.com/minio/minio-go/v7"
        "github.com/minio/minio-go/v7/pkg/credentials"
      )

      func main() {
        ${withToken ?
        `creds, _ := credentials.NewCustomTokenCredentials(\n    os.Getenv("S3_BUCKET"),\n    os.Getenv("CUSTOM_TOKEN"),\n  )` :
        `creds := credentials.NewStaticV4(\n    os.Getenv("ACCESS_KEY"),\n    os.Getenv("SECRET_KEY"),\n  )`
        }
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
  }
}
export const envExample = (bucketLink: string, siweToken?: string) =>
  siweToken ? outdent`
    # .env
    S3_BUCKET='${bucketLink}'
    CUSTOM_TOKEN='${siweToken}'
  ` : outdent`
    # .env
    S3_BUCKET='${bucketLink}'
    ACCESS_KEY='...Get your API key from the Console...'
    SECRET_KEY='...Get your API key from the Console...'
  `
