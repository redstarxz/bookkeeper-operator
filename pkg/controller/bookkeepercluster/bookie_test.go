/**
 * Copyright (c) 2018 Dell Inc., or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 */

package bookkeepercluster_test

import (
	"testing"

	"github.com/pravega/bookkeeper-operator/pkg/apis/bookkeeper/v1alpha1"
	"github.com/pravega/bookkeeper-operator/pkg/controller/bookkeepercluster"
	"github.com/pravega/bookkeeper-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBookkeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bookkeeper")
}

var _ = Describe("Bookie", func() {
	var _ = Describe("Bookie Test", func() {
		var (
			bk *v1alpha1.BookkeeperCluster
		)
		BeforeEach(func() {
			bk = &v1alpha1.BookkeeperCluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default",
				},
			}
		})
		Context("User is specifying bookkeeper journal and ledger path ", func() {

			var (
				customReq *corev1.ResourceRequirements
			)
			BeforeEach(func() {
				customReq = &corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("2"),
						corev1.ResourceMemory: resource.MustParse("4Gi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("4"),
						corev1.ResourceMemory: resource.MustParse("6Gi"),
					},
				}
				boolFalse := false
				bk.Spec = v1alpha1.BookkeeperClusterSpec{
					Version:            "0.4.0",
					ServiceAccountName: "bk-operator",
					EnvVars:            "bk-configmap",
					AutoRecovery:       &boolFalse,
					Resources:          customReq,
					Options: map[string]string{
						"journalDirectories": "/bk/journal/j0,/bk/journal/j1,/bk/journal/j2,/bk/journal/j3",
						"ledgerDirectories":  "/bk/ledgers/l0,/bk/ledgers/l1,/bk/ledgers/l2,/bk/ledgers/l3",
						"indexDirectories":   "/bk/index/i0,/bk/index/i1",
					},
				}
				bk.WithDefaults()
			})
			Context("Bookkeeper", func() {
				It("should create a headless service", func() {
					headlessservice := bookkeepercluster.MakeBookieHeadlessService(bk)
					Ω(headlessservice.Name).Should(Equal(util.HeadlessServiceNameForBookie(bk.Name)))
				})

				It("should create a pod disruption budget", func() {
					pdb := bookkeepercluster.MakeBookiePodDisruptionBudget(bk)
					Ω(pdb.Name).Should(Equal(util.PdbNameForBookie(bk.Name)))
				})

				It("should create a config-map", func() {
					cm := bookkeepercluster.MakeBookieConfigMap(bk)
					Ω(cm.Name).Should(Equal(util.ConfigMapNameForBookie(bk.Name)))
				})

				It("should create a stateful set", func() {
					ss := bookkeepercluster.MakeBookieStatefulSet(bk)
					Ω(ss.Name).Should(Equal(util.StatefulSetNameForBookie(bk.Name)))
				})

				It("should have journal and ledgers dir set to the values given by user", func() {
					sts := bookkeepercluster.MakeBookieStatefulSet(bk)
					mountledger0 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath
					mountledger1 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath
					mountledger2 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[2].MountPath
					mountledger3 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[3].MountPath
					Ω(mountledger0).Should(Equal("/bk/ledgers/l0"))
					Ω(mountledger1).Should(Equal("/bk/ledgers/l1"))
					Ω(mountledger2).Should(Equal("/bk/ledgers/l2"))
					Ω(mountledger3).Should(Equal("/bk/ledgers/l3"))
					mountjournal0 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[4].MountPath
					mountjournal1 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[5].MountPath
					mountjournal2 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[6].MountPath
					mountjournal3 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[7].MountPath
					Ω(mountjournal0).Should(Equal("/bk/journal/j0"))
					Ω(mountjournal1).Should(Equal("/bk/journal/j1"))
					Ω(mountjournal2).Should(Equal("/bk/journal/j2"))
					Ω(mountjournal3).Should(Equal("/bk/journal/j3"))
					mountindex0 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[8].MountPath
					mountindex1 := sts.Spec.Template.Spec.Containers[0].VolumeMounts[9].MountPath
					Ω(mountindex0).Should(Equal("/bk/index/i0"))
					Ω(mountindex1).Should(Equal("/bk/index/i1"))
				})
			})
		})
		Context("User is not specifying bookkeeper journal and ledger path ", func() {
			BeforeEach(func() {
				bk.Spec = v1alpha1.BookkeeperClusterSpec{}
				bk.WithDefaults()
			})
			Context("Bookkeeper", func() {
				It("should create a headless service", func() {
					headlessService := bookkeepercluster.MakeBookieHeadlessService(bk)
					Ω(headlessService.Name).Should(Equal(util.HeadlessServiceNameForBookie(bk.Name)))
				})

				It("should create a pod disruption budget", func() {
					pdb := bookkeepercluster.MakeBookiePodDisruptionBudget(bk)
					Ω(pdb.Name).Should(Equal(util.PdbNameForBookie(bk.Name)))
				})

				It("should create a config-map", func() {
					cm := bookkeepercluster.MakeBookieConfigMap(bk)
					Ω(cm.Name).Should(Equal(util.ConfigMapNameForBookie(bk.Name)))
				})

				It("should create a stateful set", func() {
					ss := bookkeepercluster.MakeBookieStatefulSet(bk)
					Ω(ss.Name).Should(Equal(util.StatefulSetNameForBookie(bk.Name)))
				})
				It("should have journal and ledgers dir set to default value", func() {
					sts := bookkeepercluster.MakeBookieStatefulSet(bk)
					mountledger := sts.Spec.Template.Spec.Containers[0].VolumeMounts[0].MountPath
					Ω(mountledger).Should(Equal("/bk/ledgers"))
					mountjournal := sts.Spec.Template.Spec.Containers[0].VolumeMounts[1].MountPath
					Ω(mountjournal).Should(Equal("/bk/journal"))
					indexjournal := sts.Spec.Template.Spec.Containers[0].VolumeMounts[2].MountPath
					Ω(indexjournal).Should(Equal("/bk/index"))
				})
			})
		})
	})
})
