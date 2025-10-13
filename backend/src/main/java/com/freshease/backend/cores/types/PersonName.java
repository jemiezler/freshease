/*
 * Click nbfs://nbhost/SystemFileSystem/Templates/Licenses/license-default.txt to change this license
 * Click nbfs://nbhost/SystemFileSystem/Templates/Classes/Class.java to edit this template
 */

package com.freshease.backend.cores.types;

import jakarta.persistence.Embeddable;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

/**
 *
 * @author jemiezler
 */
@Embeddable
@NoArgsConstructor
@Getter
@Setter
@AllArgsConstructor
@Builder
public class PersonName {
    private String first;
    private String middle;
    private String last;
}