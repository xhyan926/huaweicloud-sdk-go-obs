package com.obs.services.model.Qos;

import static org.junit.Assert.*;

import com.obs.services.model.Qos.NetworkType;
import org.junit.Test;

public class NetworkTypeTest {

    @Test
    public void testEnumValuesAndCodes() {
        // 验证枚举值的code是否正确
        assertEquals("intranet", NetworkType.INTRANET.getCode());
        assertEquals("extranet", NetworkType.EXTRANET.getCode());
        assertEquals("total", NetworkType.TOTAL.getCode());
        
        // 验证枚举值数量
        assertEquals(3, NetworkType.values().length);
    }

    @Test
    public void testGetValueFromCode_MatchCase() {
        // 测试与枚举code完全匹配的情况
        assertEquals(NetworkType.INTRANET, NetworkType.getValueFromCode("intranet"));
        assertEquals(NetworkType.EXTRANET, NetworkType.getValueFromCode("extranet"));
        assertEquals(NetworkType.TOTAL, NetworkType.getValueFromCode("total"));
    }

    @Test
    public void testGetValueFromCode_IgnoreCase() {
        // 测试大小写不匹配但内容相同的情况（方法应忽略大小写）
        assertEquals(NetworkType.INTRANET, NetworkType.getValueFromCode("INTRANET"));
        assertEquals(NetworkType.EXTRANET, NetworkType.getValueFromCode("ExTrAnEt"));
        assertEquals(NetworkType.TOTAL, NetworkType.getValueFromCode("TOTAL"));
    }

    @Test
    public void testGetValueFromCode_InvalidCode() {
        // 测试不存在的code
        assertNull(NetworkType.getValueFromCode("invalid"));
        assertNull(NetworkType.getValueFromCode(""));
        assertNull(NetworkType.getValueFromCode(null));
        assertNull(NetworkType.getValueFromCode("Intranet123"));
    }
}
    