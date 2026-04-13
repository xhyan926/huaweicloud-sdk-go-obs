package com.obs.services.model.lifecycle;

import com.obs.services.model.LifecycleConfiguration;
import org.junit.Test;
import static org.junit.Assert.*;

public class LifeCycleConfiguraionUnitTest {

    @Test
    public void testRuleEquals() {
        LifecycleConfiguration config1 = new LifecycleConfiguration();
        LifecycleConfiguration config2 = new LifecycleConfiguration();

        LifecycleConfiguration.Rule rule1 = config1.new Rule("a", "", true);
        LifecycleConfiguration.Rule rule2 = config2.new Rule("b", "", true);

        config1.addRule(rule1);
        config2.addRule(rule2);

        assertFalse(rule1.equals(rule2));

        LifecycleConfiguration.Rule rule3 = config1.new Rule("a", "", true);
        LifecycleConfiguration.Rule rule4 = config2.new Rule("a", "", true);

        assertTrue(rule3.equals(rule4));

        assertFalse(rule3.equals(null));

        assertFalse(rule3.equals(config1));
    }

    @Test
    public void testLifecycleConfigurationEquals() {
        LifecycleConfiguration configA = new LifecycleConfiguration();
        LifecycleConfiguration configB = new LifecycleConfiguration();

        LifecycleConfiguration.Rule ruleA1 = configA.new Rule("a", "", true);
        LifecycleConfiguration.Rule ruleA2 = configA.new Rule("b", "", false);
        configA.addRule(ruleA1);
        configA.addRule(ruleA2);

        LifecycleConfiguration.Rule ruleB1 = configB.new Rule("a", "", true);
        LifecycleConfiguration.Rule ruleB2 = configB.new Rule("b", "", false);
        configB.addRule(ruleB1);
        configB.addRule(ruleB2);

        assertTrue(configA.equals(configB));

        ruleB2.setEnabled(true);

        assertFalse(configA.equals(configB));

        assertFalse(configA.equals(null));

        assertFalse(configA.equals(ruleA1));
    }

    @Test
    public void testEmptyRuleList() {
        LifecycleConfiguration configC = new LifecycleConfiguration();
        LifecycleConfiguration configD = new LifecycleConfiguration();

        assertTrue(configC.equals(configD));
    }
}