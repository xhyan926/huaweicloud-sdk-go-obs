/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2024-2024. All rights reserved.
 */

package com.obs.integrated_test;

import com.obs.services.model.bpa.BucketPublicAccessBlock;

import junit.framework.TestCase;

import org.junit.Assert;

public class BucketPublicAccessBlockIT extends TestCase {
    public void testTestToString() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        String testBucketPublicAccessBlockToString = testBucketPublicAccessBlock.toString();
        Assert.assertNotNull(testBucketPublicAccessBlockToString);
        String BucketPublicAccessBlockClassName = BucketPublicAccessBlock.class.getName();
        BucketPublicAccessBlockClassName =
            BucketPublicAccessBlockClassName.substring(BucketPublicAccessBlockClassName.lastIndexOf('.') + 1);
        Assert.assertTrue(testBucketPublicAccessBlockToString.contains(BucketPublicAccessBlockClassName));
    }

    public void testGetBlockPublicACLs() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Assert.assertNull(testBucketPublicAccessBlock.getBlockPublicACLs());
    }

    public void testSetBlockPublicACLs() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Boolean testBool = false;
        testBucketPublicAccessBlock.setBlockPublicACLs(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getBlockPublicACLs());
        testBool = true;
        testBucketPublicAccessBlock.setBlockPublicACLs(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getBlockPublicACLs());
    }

    public void testGetIgnorePublicACLs() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Assert.assertNull(testBucketPublicAccessBlock.getIgnorePublicACLs());
    }

    public void testSetIgnorePublicACLs() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Boolean testBool = false;
        testBucketPublicAccessBlock.setIgnorePublicACLs(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getIgnorePublicACLs());
        testBool = true;
        testBucketPublicAccessBlock.setIgnorePublicACLs(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getIgnorePublicACLs());
    }

    public void testGetBlockPublicPolicy() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Assert.assertNull(testBucketPublicAccessBlock.getBlockPublicPolicy());
    }

    public void testSetBlockPublicPolicy() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Boolean testBool = false;
        testBucketPublicAccessBlock.setBlockPublicPolicy(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getBlockPublicPolicy());
        testBool = true;
        testBucketPublicAccessBlock.setBlockPublicPolicy(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getBlockPublicPolicy());
    }

    public void testGetRestrictPublicBuckets() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Assert.assertNull(testBucketPublicAccessBlock.getRestrictPublicBuckets());
    }

    public void testSetRestrictPublicBuckets() {
        BucketPublicAccessBlock testBucketPublicAccessBlock = new BucketPublicAccessBlock();
        Boolean testBool = false;
        testBucketPublicAccessBlock.setRestrictPublicBuckets(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getRestrictPublicBuckets());
        testBool = true;
        testBucketPublicAccessBlock.setRestrictPublicBuckets(testBool);
        assertEquals(testBool, testBucketPublicAccessBlock.getRestrictPublicBuckets());
    }
}