package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.QpsLimitConfiguration;
import org.junit.Test;

public class QpsLimitConfigurationTest {

    // 测试构造函数是否正确初始化所有字段
    @Test
    public void testConstructor() {
        // 测试正数值
        long get = 100;
        long putPostDelete = 200;
        long list = 300;
        long total = 600;

        QpsLimitConfiguration config = new QpsLimitConfiguration(get, putPostDelete, list, total);

        assertEquals(get, config.getQpsGetLimit());
        assertEquals(putPostDelete, config.getQpsPutPostDeleteLimit());
        assertEquals(list, config.getQpsListLimit());
        assertEquals(total, config.getQpsTotalLimit());

        // 测试零值
        QpsLimitConfiguration zeroConfig = new QpsLimitConfiguration(0, 0, 0, 0);
        assertEquals(0, zeroConfig.getQpsGetLimit());
        assertEquals(0, zeroConfig.getQpsPutPostDeleteLimit());
        assertEquals(0, zeroConfig.getQpsListLimit());
        assertEquals(0, zeroConfig.getQpsTotalLimit());
    }

    // 测试构造函数负值情况
    @Test(expected = IllegalArgumentException.class)
    public void testConstructorWithNegativeGet() {
        new QpsLimitConfiguration(-10, 20, 30, 60);
    }

    @Test(expected = IllegalArgumentException.class)
    public void testConstructorWithNegativePutPostDelete() {
        new QpsLimitConfiguration(10, -20, 30, 60);
    }

    @Test(expected = IllegalArgumentException.class)
    public void testConstructorWithNegativeList() {
        new QpsLimitConfiguration(10, 20, -30, 60);
    }

    @Test(expected = IllegalArgumentException.class)
    public void testConstructorWithNegativeTotal() {
        new QpsLimitConfiguration(10, 20, 30, -60);
    }

    // 测试setGet方法
    @Test
    public void testSetQpsGetLimit() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);

        // 设置正数值
        config.setQpsGetLimit(150);
        assertEquals(150, config.getQpsGetLimit());

        // 设置零值
        config.setQpsGetLimit(0);
        assertEquals(0, config.getQpsGetLimit());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testSetQpsGetLimitWithNegative() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);
        config.setQpsGetLimit(-50);
    }

    // 测试setPutPostDelete方法
    @Test
    public void testSetQpsPutPostDeleteLimit() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);

        // 设置正数值
        config.setQpsPutPostDeleteLimit(250);
        assertEquals(250, config.getQpsPutPostDeleteLimit());

        // 设置零值
        config.setQpsPutPostDeleteLimit(0);
        assertEquals(0, config.getQpsPutPostDeleteLimit());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testSetQpsPutPostDeleteLimitWithNegative() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);
        config.setQpsPutPostDeleteLimit(-80);
    }

    // 测试setList方法
    @Test
    public void testSetQpsListLimit() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);

        // 设置正数值
        config.setQpsListLimit(350);
        assertEquals(350, config.getQpsListLimit());

        // 设置零值
        config.setQpsListLimit(0);
        assertEquals(0, config.getQpsListLimit());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testSetQpsListLimitWithNegative() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);
        config.setQpsListLimit(-120);
    }

    // 测试setTotal方法
    @Test
    public void testSetQpsTotalLimit() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);

        // 设置正数值
        config.setQpsTotalLimit(750);
        assertEquals(750, config.getQpsTotalLimit());

        // 设置零值
        config.setQpsTotalLimit(0);
        assertEquals(0, config.getQpsTotalLimit());
    }

    @Test(expected = IllegalArgumentException.class)
    public void testSetQpsTotalLimitWithNegative() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(0, 0, 0, 0);
        config.setQpsTotalLimit(-300);
    }

    // 测试多次修改值的情况
    @Test
    public void testMultipleUpdates() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(10, 20, 30, 60);

        // 第一次修改
        config.setQpsGetLimit(100);
        config.setQpsPutPostDeleteLimit(200);
        config.setQpsListLimit(300);
        config.setQpsTotalLimit(600);

        assertEquals(100, config.getQpsGetLimit());
        assertEquals(200, config.getQpsPutPostDeleteLimit());
        assertEquals(300, config.getQpsListLimit());
        assertEquals(600, config.getQpsTotalLimit());

        // 第二次修改
        config.setQpsGetLimit(1000);
        config.setQpsPutPostDeleteLimit(2000);
        config.setQpsListLimit(3000);
        config.setQpsTotalLimit(6000);

        assertEquals(1000, config.getQpsGetLimit());
        assertEquals(2000, config.getQpsPutPostDeleteLimit());
        assertEquals(3000, config.getQpsListLimit());
        assertEquals(6000, config.getQpsTotalLimit());
    }

    // 测试设置负值后原值不变
    @Test(expected = IllegalArgumentException.class)
    public void testNegativeValueDoesNotChangeState() {
        QpsLimitConfiguration config = new QpsLimitConfiguration(1000, 2000, 3000, 6000);
        try {
            config.setQpsGetLimit(-100);
        } finally {
            // 验证值未改变
            assertEquals(1000, config.getQpsGetLimit());
            assertEquals(2000, config.getQpsPutPostDeleteLimit());
            assertEquals(3000, config.getQpsListLimit());
            assertEquals(6000, config.getQpsTotalLimit());
        }
    }
}