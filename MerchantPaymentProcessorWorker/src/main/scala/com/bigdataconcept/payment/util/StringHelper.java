package com.bigdataconcept.payment.util;

import java.io.UnsupportedEncodingException;
import java.net.URLDecoder;
import java.net.URLEncoder;
import java.util.logging.Level;

import static java.util.logging.Logger.getLogger;

/**
 *
 *
 */
public class StringHelper {

    /**
     * This method encodes the string.
     * @param str the string that will be encoded.
     * @return the encoded string.
     */
    public static String encode(String str){

        try {
            return URLEncoder.encode(str, "UTF-8");
        } catch (UnsupportedEncodingException ex) {
            getLogger(StringHelper.class.getName()).log(Level.SEVERE, null, ex);
            return null;
        }
    }

    /**
     * This method decodes the string.
     * @param str the string that will be decoded.
     * @return the decoded string.
     */
    public static String decode(String str){

        try {
            return URLDecoder.decode(str, "UTF-8");
        } catch (UnsupportedEncodingException ex) {
            getLogger(StringHelper.class.getName()).log(Level.SEVERE, null, ex);
            return null;
        }
    }

}