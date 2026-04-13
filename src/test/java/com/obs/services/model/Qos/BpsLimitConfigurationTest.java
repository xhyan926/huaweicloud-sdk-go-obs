package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.BpsLimitConfiguration;
import org.junit.Test;

public class BpsLimitConfigurationTest {

    @Test
    public void testConstructor() {
        // 测试正数值
        long get = 1024;
        long putPost = 2048;
        long total = 3072;

        BpsLimitConfiguration config = new BpsLimitConfiguration(get, putPost, total);

        assertEquals(get, config.getBpsGetLimit());
        assertEquals(putPost, config.getBpsPutPostLimit());
        assertEquals(total, config.getBpsTotalLimit());

        // 测试零值
        BpsLimitConfiguration zeroConfig = new BpsLimitConfiguration(0, 0, 0);
        assertEquals(0, zeroConfig.getBpsGetLimit());
        assertEquals(0, zeroConfig.getBpsPutPostLimit());
        assertEquals(0, zeroConfig.getBpsTotalLimit());
    }

    // 测试构造函数负值情况
    @Test(expected = IllegalArgumentException.class)
    public void testConstructorWithNegativeGet() {
        new BpsLimitConfiguration(-512, 1024, 1536);
    }

    @Test(expected = IllegalArgumentException.class)
    public void testConstructorWithNegativePutPost() {
        new BpsLimitConfiguration(512, -1024, 1536);
    }

    @Test(expected = IllegalArgumentException.class)
    public void testConstructorWithNegativeTotal() {
        new BpsLimitConfiguration(512, 1024, -1536);
    }

    // 测试get字段的setter和getter
    @Test
    public void testSetAndGetBpsGetLimit() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(0, 0, 0);

        config.setBpsGetLimit(5120);
        assertEquals(5120, config.getBpsGetLimit());

        config.setBpsGetLimit(0);
        assertEquals(0, config.getBpsGetLimit());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testSetBpsGetLimitWithNegative() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(0, 0, 0);
        config.setBpsGetLimit(-2048);
    }

    // 测试putPost字段的setter和getter
    @Test
    public void testSetAndGetBpsPutPostLimit() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(0, 0, 0);

        config.setBpsPutPostLimit(10240);
        assertEquals(10240, config.getBpsPutPostLimit());

        config.setBpsPutPostLimit(0);
        assertEquals(0, config.getBpsPutPostLimit());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testSetBpsPutPostLimitWithNegative() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(0, 0, 0);
        config.setBpsPutPostLimit(-4096);
    }

    // 测试total字段的setter和getter
    @Test
    public void testSetAndGetBpsTotalLimit() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(0, 0, 0);

        config.setBpsTotalLimit(20480);
        assertEquals(20480, config.getBpsTotalLimit());

        config.setBpsTotalLimit(0);
        assertEquals(0, config.getBpsTotalLimit());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testSetBpsTotalLimitWithNegative() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(0, 0, 0);
        config.setBpsTotalLimit(-8192);
    }

    // 测试多次更新字段值的场景
    @Test
    public void testMultipleUpdates() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(100, 200, 300);

        // 第一次更新
        config.setBpsGetLimit(1000);
        config.setBpsPutPostLimit(2000);
        config.setBpsTotalLimit(3000);

        assertEquals(1000, config.getBpsGetLimit());
        assertEquals(2000, config.getBpsPutPostLimit());
        assertEquals(3000, config.getBpsTotalLimit());

        // 第二次更新
        config.setBpsGetLimit(5000);
        config.setBpsPutPostLimit(6000);
        config.setBpsTotalLimit(11000);

        assertEquals(5000, config.getBpsGetLimit());
        assertEquals(6000, config.getBpsPutPostLimit());
        assertEquals(11000, config.getBpsTotalLimit());
    }

    // 测试设置负值后原值不变
    @Test(expected = IllegalArgumentException.class)
    public void testNegativeValueDoesNotChangeState() {
        BpsLimitConfiguration config = new BpsLimitConfiguration(5000, 6000, 11000);
        try {
            config.setBpsGetLimit(-100);
        } finally {
            // 验证值未改变
            assertEquals(5000, config.getBpsGetLimit());
            assertEquals(6000, config.getBpsPutPostLimit());
            assertEquals(11000, config.getBpsTotalLimit());
        }
    }
}