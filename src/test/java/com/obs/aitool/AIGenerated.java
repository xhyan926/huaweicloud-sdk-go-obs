package com.obs.aitool;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

@Retention(RetentionPolicy.RUNTIME)
@Target(ElementType.METHOD)
public @interface AIGenerated {

    /**
     * The author of the generated code.
     *
     * @return the author name
     */
    String author();

    /**
     * The date when the code was generated.
     * Should follow YYYY-MM-DD format.
     *
     * @return the generation date in YYYY-MM-DD format
     */
    String date();


    /**
     * Description of the test scenario being tested.
     * Provides context about what specific behavior or condition the test verifies.
     *
     * @return description of the test scenario
     */
    String description() default "";
}
