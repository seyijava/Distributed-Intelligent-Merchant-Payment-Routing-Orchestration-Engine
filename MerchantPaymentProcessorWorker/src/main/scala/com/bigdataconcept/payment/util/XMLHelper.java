package com.bigdataconcept.payment.util;


import org.json.JSONObject;
import org.json.XML;

/**
 *
 *
 */
public class XMLHelper {

    /**
     * This method converts xml into JSONObject.
     *
     * @param xml the xml string than will be parsed.
     * @return the JSONObject json representation of xml
     */
    public static JSONObject toJson(String xml) {

        return XML.toJSONObject(xml);
    }

    /**
     * This method converts JSONObject to xml
     *
     * @param object the JSONObject that will be parsed.
     * @return the xml representation of JSONObject.
     */
    public static String fromJson(JSONObject object) {

        return XML.toString(object);
    }

}