// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package s3_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/service/s3"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfs3 "github.com/hashicorp/terraform-provider-aws/internal/service/s3"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccS3BucketPublicAccessBlock_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var config types.PublicAccessBlockConfiguration
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_public_access_block.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketPublicAccessBlockDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketExists(ctx, "aws_s3_bucket.test"),
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "bucket", rName),
					resource.TestCheckResourceAttr(resourceName, "block_public_acls", "false"),
					resource.TestCheckResourceAttr(resourceName, "block_public_policy", "false"),
					resource.TestCheckResourceAttr(resourceName, "ignore_public_acls", "false"),
					resource.TestCheckResourceAttr(resourceName, "restrict_public_buckets", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccS3BucketPublicAccessBlock_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var config types.PublicAccessBlockConfiguration
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_public_access_block.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketPublicAccessBlockDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfs3.ResourceBucketPublicAccessBlock(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3BucketPublicAccessBlock_Disappears_bucket(t *testing.T) {
	ctx := acctest.Context(t)
	var config types.PublicAccessBlockConfiguration
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_public_access_block.test"
	bucketResourceName := "aws_s3_bucket.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketPublicAccessBlockDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfs3.ResourceBucket(), bucketResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccS3BucketPublicAccessBlock_blockPublicACLs(t *testing.T) {
	ctx := acctest.Context(t)
	var config types.PublicAccessBlockConfiguration
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_public_access_block.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketPublicAccessBlockDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, true, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "block_public_acls", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "block_public_acls", "false"),
				),
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, true, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "block_public_acls", "true"),
				),
			},
		},
	})
}

func TestAccS3BucketPublicAccessBlock_blockPublicPolicy(t *testing.T) {
	ctx := acctest.Context(t)
	var config types.PublicAccessBlockConfiguration
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_public_access_block.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketPublicAccessBlockDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "block_public_policy", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "block_public_policy", "false"),
				),
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, true, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "block_public_policy", "true"),
				),
			},
		},
	})
}

func TestAccS3BucketPublicAccessBlock_ignorePublicACLs(t *testing.T) {
	ctx := acctest.Context(t)
	var config types.PublicAccessBlockConfiguration
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_public_access_block.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketPublicAccessBlockDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "ignore_public_acls", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "ignore_public_acls", "false"),
				),
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "ignore_public_acls", "true"),
				),
			},
		},
	})
}

func TestAccS3BucketPublicAccessBlock_restrictPublicBuckets(t *testing.T) {
	ctx := acctest.Context(t)
	var config types.PublicAccessBlockConfiguration
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_s3_bucket_public_access_block.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.S3EndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketPublicAccessBlockDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "restrict_public_buckets", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "restrict_public_buckets", "false"),
				),
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName, false, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBucketPublicAccessBlockExists(ctx, resourceName, &config),
					resource.TestCheckResourceAttr(resourceName, "restrict_public_buckets", "true"),
				),
			},
		},
	})
}

func TestAccS3BucketPublicAccessBlock_directoryBucket(t *testing.T) {
	ctx := acctest.Context(t)
	name := fmt.Sprintf("tf-test-bucket-%d", sdkacctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, s3.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckBucketDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketPublicAccessBlockConfig_directoryBucket(name, "false", "false", "false", "false"),
				ExpectError: regexache.MustCompile(`NoSuchBucket`), // Waiting for resource migration to AWS SDK for Go v2.
			},
		},
	})
}

func testAccCheckBucketPublicAccessBlockDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Client(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_s3_bucket_public_access_block" {
				continue
			}

			_, err := tfs3.FindPublicAccessBlockConfiguration(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("S3 Bucket Public Access Block %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckBucketPublicAccessBlockExists(ctx context.Context, n string, v *types.PublicAccessBlockConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).S3Client(ctx)

		output, err := tfs3.FindPublicAccessBlockConfiguration(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccBucketPublicAccessBlockConfig_basic(bucketName string, blockPublicAcls, blockPublicPolicy, ignorePublicAcls, restrictPublicBuckets bool) string {
	return fmt.Sprintf(`
resource "aws_s3_bucket" "test" {
  bucket = %[1]q
}

resource "aws_s3_bucket_public_access_block" "test" {
  bucket = aws_s3_bucket.test.bucket

  block_public_acls       = %[2]t
  block_public_policy     = %[3]t
  ignore_public_acls      = %[4]t
  restrict_public_buckets = %[5]t
}
`, bucketName, blockPublicAcls, blockPublicPolicy, ignorePublicAcls, restrictPublicBuckets)
}

func testAccBucketPublicAccessBlockConfig_directoryBucket(bucketName, blockPublicAcls, blockPublicPolicy, ignorePublicAcls, restrictPublicBuckets string) string {
	return acctest.ConfigCompose(testAccDirectoryBucketConfig_base(bucketName), fmt.Sprintf(`
resource "aws_s3_directory_bucket" "test" {
  bucket = local.bucket
  location {
    name = local.location_name
  }
}
resource "aws_s3_bucket_public_access_block" "bucket" {
  bucket = aws_s3_directory_bucket.test.bucket
  block_public_acls       = %[2]q
  block_public_policy     = %[3]q
  ignore_public_acls      = %[4]q
  restrict_public_buckets = %[5]q
}
`, bucketName, blockPublicAcls, blockPublicPolicy, ignorePublicAcls, restrictPublicBuckets))
}
