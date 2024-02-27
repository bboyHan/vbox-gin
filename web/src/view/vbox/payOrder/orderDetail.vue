<template>
  <div>
    <div v-show="dialogCountVisible">
      <div class="c_container">
        <div class="c_content">
          <count-down
              :fire="fire"
              :tiping="tiping"
              :tipend="tipend"
              time="12"
              @statusChange="onStatusChange"
              @end="onEnd"
              :statusChange="[2000,500]"
              width="180"
              height="180"
          >
          </count-down>
          <div class="buttons">
            <el-row :gutter="6">
              <el-col :span="24">
                <el-button type="success" class="c_button" round>正在通过安全验证，请等待...</el-button>
              </el-col>
              <el-col :span="24">
                <el-button type="primary" class="c_button" round>订单正在匹配中，预计5-20秒</el-button>
              </el-col>
            </el-row>
          </div>
        </div>
      </div>
    </div>

    <div v-show="payVisible">
      <!-- 显示新的 div 的代码... -->
      <!--   1000 引导   -->
      <div v-if="payTypeVisible >= 1000 && payTypeVisible < 1099">
        <div class="p_container">
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content" :style="backgroundImageStyle">
            <el-row :gutter="12">
              <el-col style="width: 80px; height: 80px">
              </el-col>
              <el-col>
                <!--                <div style="color: #6B7687; margin-top: 10px; font-size: 16px">无法充值或提示错误，请联系客服！</div>-->
              </el-col>
              <el-col>
                <p v-if="Number(payData.channel_code) === 1003" style="color: red;margin-right: 20px;margin-left: 20px">
                  充值时必须<b style="color: blue;">滑动选择【王者荣耀】并粘贴账号</b>付款即可到账！</p>
                <p v-else-if="Number(payData.channel_code) === 1006"
                   style="color: red;margin-right: 20px;margin-left: 20px">
                  点击<b style="color: blue;">【一键复制】</b>跳转微信，打开任意<b style="color: blue;">【聊天窗口】</b>粘贴发送复制内容，再根据提示步骤操作支付！</p>
                <p v-else style="color: red;margin-right: 20px;margin-left: 20px">充值前<b style="color: blue;">核对【订单金额】并复制账号</b>，根据指导步骤付款即可到账！</p>
              </el-col>
              <el-col>
                <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
              </el-col>
              <el-col>
                <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="24">

              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_inner" :style="backgroundImageStyle" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                </div>
                <div class="medicine-bag">
                  <span>{{ payData.account }}</span>
                </div>
                <div v-if="Number(payData.channel_code) === 1006" class="copy-tip">
                  <span>点复制</span>
                  <span class="jtone"></span>
                  <span>跳转微信</span>
                  <span class="jttwo"></span>
                  <span>聊天粘贴</span>
                  <span class="jtthree"></span>
                  <span>提示支付</span>
                </div>
                <div v-else class="copy-tip">
                  <span>长按框内</span>
                  <span class="jtone"></span>
                  <span>复制</span>
                  <span class="jttwo"></span>
                  <span>记金额</span>
                  <span class="jtthree"></span>
                  <span>打开跳转</span>
                </div>
                <div v-if="!copyInfoVisible">
                  <button class="btn-copy copy_button" @click="copyInfo">① 一键复制</button>
                </div>
                <div v-else>
                  <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button">
          <el-row :gutter="12">
            <el-col :span="24">
              <div v-if="Number(payData.channel_code) === 1003">
                <button class="btn-copy p_button" @click="openYdVisible" style="font-size: 15px">②
                  点击付款,右滑选"王者荣耀"
                </button>
              </div>
              <div v-else>
                <button class="btn-copy p_button" @click="openYdVisible">② 点击付款</button>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   1100 JW 引导   -->
      <div v-if="payTypeVisible >= 1100 && payTypeVisible < 1199">
        <div class="p_container">
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content" :style="backgroundImageStyle">
            <el-row :gutter="12">
              <el-col style="width: 80px; height: 80px">
              </el-col>
              <!--              <el-col>
                              <img src="@/assets/header.png" alt="" style="width: 80px; height: 80px">
                            </el-col>-->
              <el-col>
                <div style="color: #6B7687; margin-top: 10px; font-size: 16px">无法充值或提示错误，请联系客服！</div>
              </el-col>
              <el-col>
                <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
              </el-col>
              <el-col>
                <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="24">

              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_card_info_button">
          <el-row :gutter="12">
            <el-col :span="24">
              <button class="btn-copy p_button" @click="openYdVisible">① 点击付款</button>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_card_info_inner" :style="backgroundImageStyle" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-jw-bag">
                  <textarea v-model="inputString"
                            :placeholder="`粘贴示例：\n您已购买成功(订单号:205...)，如下：\n卡号：2312290766321121;\n密码：2732221581323347;`"></textarea>
                  <!--                  <el-input v-model="inputString" placeholder="请输入待匹配的字符串"-->
                  <!--                            style=" border: none;background-color: transparent;"></el-input>-->
                  <!-- 在这里显示匹配到的卡号和密码 -->
                </div>
                <div>
                  <button class="btn-copy copy_button" @click="">② 粘贴智能识别</button>
                  <div class="medicine-jw-card-info-bag">
                    <div v-if="card1100Number && password1100" class="result-container">
                      <p><strong style="padding-right: 10px;font-size: 14px">卡号:</strong><b
                          style="color: blue;font-size: 16px"> {{ card1100Number }}</b></p>
                      <p><strong style="padding-right: 10px;font-size: 14px">密码:</strong><b
                          style="color: blue;font-size: 16px"> {{ password1100 }}</b></p>
                    </div>
                    <div v-else>
                      <p style="color: #c4bdbd;">未识别到卡号和密码，请核对是否包含16位卡号和16位密码</p>
                    </div>
                  </div>
                </div>
                <div class="p_content_card_submit_button">
                  <div v-if="card1100Number && password1100" class="result-container">
                    <el-row :gutter="12">
                      <el-col :span="24">
                        <button class="btn-copy p_submit_success_button" @click="openCard1100Visible">③ 提交卡密</button>
                      </el-col>
                    </el-row>
                  </div>
                  <div v-else>
                    <button class="btn-copy p_submit_button" @click="warnCardInfo">③ 提交卡密</button>
                  </div>
                </div>

              </div>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   1200 引导   -->
      <div v-if="payTypeVisible >= 1200 && payTypeVisible < 1299">
        <div class="p_container">
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content" :style="backgroundImageStyle">
            <el-row :gutter="12">
              <el-col style="width: 80px; height: 80px">
              </el-col>
              <!--              <el-col>
                              <img src="@/assets/header.png" alt="" style="width: 80px; height: 80px">
                            </el-col>-->
              <el-col>
                <!--                <div style="color: #6B7687; margin-top: 10px; font-size: 16px">无法充值或提示错误，请联系客服！</div>-->
              </el-col>
              <el-col>
                <p style="color: red;margin-right: 20px;margin-left: 20px">充值前<b style="color: blue;">核对【订单金额】并复制账号</b>，根据指导步骤付款即可到账！</p>
              </el-col>
              <el-col>
                <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
              </el-col>
              <el-col>
                <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="24">
              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_inner" :style="backgroundImageStyle" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                </div>
                <div class="medicine-bag">
                  <span>{{ payData.account }}</span>
                </div>
                <div class="copy-tip">
                  <span>长按框内</span>
                  <span class="jtone"></span>
                  <span>复制</span>
                  <span class="jttwo"></span>
                  <span>记金额</span>
                  <span class="jtthree"></span>
                  <span>打开跳转</span>
                </div>
                <div v-if="!copyInfoVisible">
                  <button class="btn-copy copy_button" @click="copyInfo">① 一键复制</button>
                </div>
                <div v-else>
                  <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button">
          <el-row :gutter="12">
            <el-col :span="24">
              <div>
                <button class="btn-copy p_button" @click="openYdVisible">② 点击付款</button>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   2000 引导   -->
      <div v-if="payTypeVisible >= 2000 && payTypeVisible < 2099">
        <div class="p_container">
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content" :style="backgroundImageStyle">
            <el-row :gutter="12">
              <el-col style="width: 80px; height: 80px">
              </el-col>
              <!--              <el-col>
                              <img src="@/assets/header.png" alt="" style="width: 80px; height: 80px">
                            </el-col>-->
              <el-col>
                <p style="color: red;margin-right: 20px;margin-left: 20px">充值前<b style="color: blue;">核对【订单金额】并复制账号</b>，根据指导步骤付款即可到账！</p>
              </el-col>
              <el-col>
                <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
              </el-col>
              <el-col>
                <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="24">
              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_inner" :style="backgroundImageStyle" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                </div>
                <div class="medicine-bag">
                  <span>{{ payData.account }}</span>
                </div>
                <div class="copy-tip">
                  <span>长按框内</span>
                  <span class="jtone"></span>
                  <span>复制</span>
                  <span class="jttwo"></span>
                  <span>记金额</span>
                  <span class="jtthree"></span>
                  <span>打开跳转</span>
                </div>
                <div v-if="!copyInfoVisible">
                  <button class="btn-copy copy_button" @click="copyInfo">① 一键复制</button>
                </div>
                <div v-else>
                  <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button">
          <el-row :gutter="12">
            <el-col :span="24">
              <button class="btn-copy p_button" @click="openYdVisible">② 点击付款</button>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   3000 直付   -->
      <div v-if="payTypeVisible >= 3000 && payTypeVisible < 3099">
        <div class="p_container">
          <!--        <div class="p_blue-section" v-for="(color, index) in blueColors" :key="index" :style="{ backgroundColor: color }"></div>-->
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content">
            <el-row :gutter="12">
              <el-col style="width: 80px; height: 80px">
              </el-col>
              <!--              <el-col>
                              <img src="@/assets/header.png" alt="" style="width: 80px; height: 80px">
                            </el-col>-->
              <el-col>
                <div style="color: #6B7687; margin-top: 10px; font-size: 16px">无法充值或提示错误，请联系客服！</div>
              </el-col>
              <el-col>
                <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
              </el-col>
              <el-col>
                <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col>
                <div style="height: 100px;">
                </div>
              </el-col>
              <el-col :span="24">
                <!--                <el-button class="p_button" type="primary">复制充值账号 立即付款</el-button>-->
              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_inner" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 5px; margin-top: 5px">
              </div>
            </el-col>
            <el-col>
              <div>
                <div class="medicine-bag-qr">
                  <div v-if="qrcodeUrl">
                    <img :src="qrcodeUrl" alt="QR Code" style="height: 174px"/>
                  </div>
                  <div v-else>
                    <span>暂无二维码</span>
                  </div>
                </div>
                <div class="copy-tip">
                  <span>切换手机</span>
                  <span class="jtone"></span>
                  <span>打开微信</span>
                  <span class="jttwo"></span>
                  <span>扫一扫</span>
                  <span class="jtthree"></span>
                  <span>扫码付款</span>
                </div>
                <!--                <button class="copy_button" @click="">一键复制</button>-->
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button_qr">
          <el-row :gutter="12">
            <el-col>
              <button class="p_button" @click="">扫码直付</button>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   4000 引导   -->
      <div v-if="payTypeVisible >= 4000 && payTypeVisible < 4099">
        <div class="p_container">
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content" :style="backgroundImageStyle">
            <el-row :gutter="12">
              <el-col style="width: 80px; height: 80px">
              </el-col>
              <!--              <el-col>
                              <img src="@/assets/header.png" alt="" style="width: 80px; height: 80px">
                            </el-col>-->
              <el-col>
                <!--                <div style="color: #6B7687; margin-top: 10px; font-size: 16px">无法充值或提示错误，请联系客服！</div>-->
              </el-col>
              <el-col>
                <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
              </el-col>
              <el-col>
                <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="24">
              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_inner" :style="backgroundImageStyle" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                </div>
                <div class="medicine-bag">
                  <span>{{ payData.account }}</span>
                </div>
                <div class="copy-tip">
                  <span>长按框内</span>
                  <span class="jtone"></span>
                  <span>复制</span>
                  <span class="jttwo"></span>
                  <span>记金额</span>
                  <span class="jtthree"></span>
                  <span>打开跳转</span>
                </div>
                <div v-if="!copyInfoVisible">
                  <button class="btn-copy copy_button" @click="copyInfo">① 一键复制</button>
                </div>
                <div v-else>
                  <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_button">
          <el-row :gutter="12">
            <el-col :span="24">
              <button class="btn-copy p_button" @click="openYdVisible">② 点击付款</button>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   6000 卡密 引导   -->
      <div v-if="payTypeVisible >= 6000 && payTypeVisible < 6099">
        <div class="p_container">
          <div class="p_blue-section" v-for="index in 10" :key="index"
               :style="{ backgroundColor: generateColor(index) }"></div>
          <div class="p_content" :style="backgroundImageStyle">
            <el-row :gutter="12">
              <el-col style="width: 60px; height: 60px;text-align: center; font-size: 20px;margin-top: 20px;margin-bottom:-20px;color: #6B7687;">充值须知
              </el-col>
              <el-col style="text-align: left">
                <div style="color: red;margin-right: 20px;margin-left: 20px;">1. 充值前<b style="color: blue;">核对【订单金额】</b></div>
                <div style="color: red;margin-right: 20px;margin-left: 20px;">2. 可自行前往<b style="color: blue;">京东/淘宝/抖音/各大商城</b>购买卡密</div>
                <div style="color: red;margin-right: 20px;margin-left: 20px;">3. 根据指导步骤<b style="color: blue;">付款并获取卡号</b></div>
                <div style="color: red;margin-right: 20px;margin-left: 20px;">4.<b style="color: blue;"> 复制卡号</b>在下方框输入进行<b style="color: blue;">提交</b></div>
              </el-col>
              <el-col>
                <div style="color: #6B7687; margin-top: 20px; font-size: 60px">￥{{ payData.money }}.00</div>
              </el-col>
              <el-col>
                <div style="color: #e81239; margin-top: 10px; font-size: 16px">
                  <el-icon style="margin-right: 5px">
                    <WarningFilled/>
                  </el-icon>
                  请在规定时间内付款！
                  <div>
                    <span v-if="countdowns[0] > 0">{{ formatTime(countdowns[0]) }} </span>
                    <span v-else>-1 （已过期）</span>
                  </div>
                </div>
              </el-col>
              <el-col :span="24">

              </el-col>
            </el-row>
          </div>
        </div>
        <div class="p_content_card_info_button">
          <el-row :gutter="12">
            <el-col :span="24">
              <button class="btn-copy p_button" @click="openYdVisible" style="font-size: 16px">① 点我购买"京东E卡"</button>
            </el-col>
          </el-row>
        </div>
        <div class="p_content_card_info_inner" :style="backgroundImageStyle" style="margin-top: 20px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-ec-bag">
                  <textarea v-model="inputString"
                            :placeholder="`粘贴示例（卡号）：E71D-5E47-33ED-50BA`"></textarea>
                  <!--                  <el-input v-model="inputString" placeholder="请输入待匹配的字符串"-->
                  <!--                            style=" border: none;background-color: transparent;"></el-input>-->
                  <!-- 在这里显示匹配到的卡号和密码 -->
                </div>
                <div>
                  <button class="btn-copy copy_button" @click="">② 粘贴智能识别</button>
                  <div class="medicine-jw-card-info-bag">
                    <div v-if="card6000Number" class="result-container">
                      <p><strong style="padding-right: 10px;font-size: 14px">卡号:</strong><b
                          style="color: blue;font-size: 16px"> {{ card6000Number }}</b></p>
                    </div>
                    <div v-else>
                      <p style="color: #c4bdbd;">未识别到卡号，请核对是否包含16位卡号</p>
                    </div>
                  </div>
                </div>
                <div class="p_content_card_submit_button">
                  <div v-if="card6000Number" class="result-container">
                    <el-row :gutter="12">
                      <el-col :span="24">
                        <button class="btn-copy p_submit_success_button" @click="open6000CardVisible">③ 提交卡号</button>
                      </el-col>
                    </el-row>
                  </div>
                  <div v-else>
                    <button class="btn-copy p_submit_button" @click="warnCardInfo">③ 提交卡号</button>
                  </div>
                </div>

              </div>
            </el-col>
          </el-row>
        </div>
      </div>

      <!--   jw卡密信息确认   -->
      <el-dialog width="360px" v-model="dialog1100CardVisible" :draggable="true" :before-close="close1100CardVisible"
                 :style="backgroundYdImageStyle" top="40vh" destroy-on-close>
        <div>
          <div>
            <div class="medicine-jw-card-info-submit-bag">
              <div v-if="card1100Number && password1100" class="result-container">
                <p style="padding: 5px"><strong style="color: red">核对确认，提交后不可修改！</strong></p>
                <p style="padding: 5px"><strong>卡号:</strong><b style="color: blue;font-size: 20px"> {{
                    card1100Number
                  }}</b></p>
                <p style="padding: 5px"><strong>密码:</strong><b style="color: blue;font-size: 20px"> {{ password1100 }}</b>
                </p>
              </div>
              <div v-else>
                <p>未识别到卡号和密码，请核对是否包含16位卡号和16位密码</p>
              </div>
            </div>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <div class="yd_p_content_button_qr">
              <el-row :gutter="12">
                <el-col>
                  <div>
                    <button class="yd_p_button" @click="submit1100CardInfo">确认提交卡密</button>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </template>
      </el-dialog>

      <!--   ec卡密信息确认   -->
      <el-dialog width="360px" v-model="dialog6000CardVisible" :draggable="true" :before-close="close6000CardVisible"
                 :style="backgroundYdImageStyle" top="40vh" destroy-on-close>
        <div>
          <div>
            <div class="medicine-ec-card-info-submit-bag">
              <div v-if="card6000Number" class="result-container">
                <p style="padding: 5px"><strong style="color: red">核对确认，提交后不可修改！</strong></p>
                <p style="padding: 5px"><strong>卡号:</strong><b style="color: blue;font-size: 20px"> {{
                    card6000Number
                  }}</b></p>
              </div>
              <div v-else>
                <p>未识别到卡号，请核对是否包含16位卡号</p>
              </div>
            </div>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <div class="yd_p_content_button_qr">
              <el-row :gutter="12">
                <el-col>
                  <div>
                    <button class="yd_p_button" @click="submit6000CardInfo">确认提交卡号</button>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </template>
      </el-dialog>

      <!--   步骤指导   -->
      <!--   1000   -->
      <el-dialog width="360px" v-model="dialogYd1000Visible" :draggable="true" :before-close="closeYdDialog"
                 :style="backgroundYdImageStyle"
                 top="5vh" destroy-on-close>
        <div style="padding: 0; margin: -20px 0 0;">
          <div>
            <div v-if="Number(payData.channel_code) === 1003">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_jym.png">
            </div>
            <div v-else-if="Number(payData.channel_code) === 1002">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_dy.png">
            </div>
            <div v-else-if="Number(payData.channel_code) === 1001">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_jd.png">
            </div>
            <div v-else-if="Number(payData.channel_code) === 1004">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_zfb.png">
            </div>
            <div v-else-if="Number(payData.channel_code) === 1005">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_tb.png">
            </div>
            <div v-else-if="Number(payData.channel_code) === 1006">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_wx.png">
            </div>
            <div v-else-if="Number(payData.channel_code) === 1007">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_pdd.png">
            </div>
          </div>
          <div class="p_content_yd_inner" :style="backgroundImageStyle" style="margin-top: 10px;">
            <el-row>
              <el-col>
                <div style="height: 100px; margin-top: 20px">
                  <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                  </div>
                  <div class="medicine-bag">
                    <span>{{ payData.account }}</span>
                  </div>
                  <div v-if="Number(payData.channel_code) === 1006" class="copy-tip">
                    <span>点复制</span>
                    <span class="jtone"></span>
                    <span>跳微信</span>
                    <span class="jttwo"></span>
                    <span>聊天粘贴</span>
                    <span class="jtthree"></span>
                    <span>看步骤</span>
                  </div>
                  <div v-else class="copy-tip">
                    <span>长按框内</span>
                    <span class="jtone"></span>
                    <span>复制</span>
                    <span class="jttwo"></span>
                    <span>记金额</span>
                    <span class="jtthree"></span>
                    <span>打开跳转</span>
                  </div>

                  <div v-if="!copyInfoVisible">
                    <button class="btn-copy copy_button" @click="copyInfo">一键复制</button>
                  </div>
                  <div v-else>
                    <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <div class="yd_p_content_button_qr">
              <el-row :gutter="12">
                <el-col>
                  <div v-if="readInfoVisible">
                    <button class="yd_read_p_button" @click="">我已阅读并知晓({{ countdownTime }}s)</button>
                  </div>
                  <div v-else>
                    <div v-if="Number(payData.channel_code) === 1003">
                      <button class="btn-copy yd_p_button" @click="openPay" style="font-size: 16px">
                        点此支付,右滑选"王者荣耀"
                      </button>
                    </div>
                    <div v-else>
                      <button class="btn-copy yd_p_button" @click="openPay">点此支付</button>
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </template>
      </el-dialog>

      <!--   引导步骤1100   -->
      <el-dialog width="360px" v-model="dialogYd1100Visible" :draggable="true" :before-close="closeYdDialog"
                 :style="backgroundYdImageStyle"
                 top="5vh" destroy-on-close>
        <div style="padding: 0; margin: -20px 0 0;">
          <div>
            <div v-if="Number(payData.channel_code) === 1101">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_qb_jw.png">
            </div>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <div class="yd_p_content_button_qr">
              <el-row :gutter="12">
                <el-col>
                  <div v-if="readInfoVisible">
                    <button class="yd_read_p_button" @click="">我已阅读并知晓({{ countdownTime }}s)</button>
                  </div>
                  <div v-else>
                    <button class="btn-copy yd_p_button" @click="openPay">点此支付</button>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </template>
      </el-dialog>

      <!--   引导步骤1200   -->
      <el-dialog width="360px" v-model="dialogYd1200Visible" :draggable="true" :before-close="closeYdDialog"
                 :style="backgroundYdImageStyle"
                 top="5vh" destroy-on-close>
        <div style="padding: 0; margin: -20px 0 0;">
          <div>
            <div v-if="Number(payData.channel_code) === 1201">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_dnf_tb.png">
            </div>
            <div v-else-if="Number(payData.channel_code) === 1202">
              <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                   src="@/assets/yd_dnf_jd.png">
            </div>
          </div>
          <div class="p_content_yd_inner" :style="backgroundImageStyle" style="margin-top: 10px;">
            <el-row>
              <el-col>
                <div style="height: 100px; margin-top: 20px">
                  <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                  </div>
                  <div class="medicine-bag">
                    <span>{{ payData.account }}</span>
                  </div>
                  <div class="copy-tip">
                    <span>长按框内</span>
                    <span class="jtone"></span>
                    <span>复制</span>
                    <span class="jttwo"></span>
                    <span>记金额</span>
                    <span class="jtthree"></span>
                    <span>打开跳转</span>
                  </div>
                  <div v-if="!copyInfoVisible">
                    <button class="btn-copy copy_button" @click="copyInfo">一键复制</button>
                  </div>
                  <div v-else>
                    <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <div class="yd_p_content_button_qr">
              <el-row :gutter="12">
                <el-col>
                  <div v-if="readInfoVisible">
                    <button class="yd_read_p_button" @click="">我已阅读并知晓({{ countdownTime }}s)</button>
                  </div>
                  <div v-else>
                    <button class="btn-copy yd_p_button" @click="openPay">点此支付</button>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </template>
      </el-dialog>

      <!--   引导步骤2000   -->
      <el-dialog width="360px" v-model="dialogYd2000Visible" :draggable="true" :before-close="closeYdDialog"
                 :style="backgroundYdImageStyle"
                 top="5vh" destroy-on-close>
        <div style="padding: 0; margin: -20px 0 0;">
          <div>
            <div v-if="Number(payData.channel_code) === 2001">
              <div v-if="tmH5Visible">
                <img alt style="width: 100%; height: 100%;border-radius: 20px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                     src="@/assets/yd_j3_tm.png">
              </div>
              <div v-else>
                <img alt style="width: 100%; height: 100%;border-radius: 20px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                     src="@/assets/yd_j3_tb.png">
              </div>
            </div>
          </div>
          <div class="p_content_yd_inner" :style="backgroundImageStyle" style="margin-top: 10px;">
            <el-row>
              <el-col>
                <div style="height: 100px; margin-top: 20px">
                  <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                  </div>
                  <div class="medicine-bag">
                    <span>{{ payData.account }}</span>
                  </div>
                  <div class="copy-tip">
                    <span>长按框内</span>
                    <span class="jtone"></span>
                    <span>复制</span>
                    <span class="jttwo"></span>
                    <span>记金额</span>
                    <span class="jtthree"></span>
                    <span>打开跳转</span>
                  </div>
                  <div v-if="!copyInfoVisible">
                    <button class="btn-copy copy_button" @click="copyInfo">一键复制</button>
                  </div>
                  <div v-else>
                    <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
        <template #footer>
          <div class="dialog-footer">
            <div class="yd_p_content_button_qr">
              <el-row :gutter="12">
                <el-col>
                  <div v-if="readInfoVisible">
                    <button class="yd_read_p_button" @click="">我已阅读并知晓({{ countdownTime }}s)</button>
                  </div>
                  <div v-else>
                    <button class="btn-copy yd_p_button" @click="openPay">点此支付</button>
                  </div>
                </el-col>
              </el-row>
            </div>
          </div>
        </template>
      </el-dialog>
    </div>

    <!--   步骤指导4000   -->
    <el-dialog width="360px" v-model="dialogYd4000Visible" :draggable="true" :before-close="closeYdDialog"
               :style="backgroundYdImageStyle"
               top="5vh" destroy-on-close>
      <div style="padding: 0; margin: -20px 0 0;">
        <div>
          <div v-if="Number(payData.channel_code) === 4001">
            <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                 src="@/assets/yd_sdo_tb.png">
          </div>
        </div>
        <div class="p_content_yd_inner" :style="backgroundImageStyle" style="margin-top: 10px;">
          <el-row>
            <el-col>
              <div style="height: 100px; margin-top: 20px">
                <div class="medicine-money-bag">
                  <span><span style="color: red">牢记</span>充值金额：<span style="color: blue">￥{{
                      payData.money
                    }}.00</span></span>
                </div>
                <div class="medicine-bag">
                  <span>{{ payData.account }}</span>
                </div>
                <div class="copy-tip">
                  <span>长按框内</span>
                  <span class="jtone"></span>
                  <span>复制</span>
                  <span class="jttwo"></span>
                  <span>记金额</span>
                  <span class="jtthree"></span>
                  <span>打开跳转</span>
                </div>
                <div v-if="!copyInfoVisible">
                  <button class="btn-copy copy_button" @click="copyInfo">一键复制</button>
                </div>
                <div v-else>
                  <button class="btn-copy copy_success_button" @click="copyInfo">复制成功</button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <div class="yd_p_content_button_qr">
            <el-row :gutter="12">
              <el-col>
                <div v-if="readInfoVisible">
                  <button class="yd_read_p_button" @click="">我已阅读并知晓({{ countdownTime }}s)</button>
                </div>
                <div v-else>
                  <button class="btn-copy yd_p_button" @click="openPay">点此支付</button>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </template>
    </el-dialog>

    <!--   引导步骤6000   -->
    <el-dialog width="360px" v-model="dialogYd6000Visible" :draggable="true" :before-close="closeYdDialog"
               :style="backgroundYdImageStyle"
               top="5vh" destroy-on-close>
      <div style="padding: 0; margin: -20px 0 0;">
        <el-carousel :interval="4000" height="500px">
          <el-carousel-item>
            <div>
              <div v-if="Number(payData.channel_code) === 6001">
                <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                     src="@/assets/yd_ec_card.jpg">
              </div>
            </div>
          </el-carousel-item>
<!--          <el-carousel-item>
            <div>
              <div v-if="Number(payData.channel_code) === 6001">
                <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                     src="@/assets/yd_ec_card.png">
              </div>
            </div>
          </el-carousel-item>-->
        </el-carousel>
<!--        <div>
          <div v-if="Number(payData.channel_code) === 6001">
            <img alt style="width: 100%; height: 100%;border-radius: 5px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);"
                 src="@/assets/yd_ec_card.png">
          </div>
        </div>-->
      </div>
      <template #footer>
        <div class="dialog-footer">
          <div class="yd_p_content_button_qr">
            <el-row :gutter="12">
              <el-col>
                <div v-if="readInfoVisible">
                  <button class="yd_read_p_button" @click="">我已阅读并知晓({{ countdownTime }}s)</button>
                </div>
                <div v-else>
                  <button class="btn-copy yd_p_jd_button" @click="openPayHref('jd')">前往京东支付</button>
                  <button class="btn-copy yd_p_tb_button" @click="openPayHref('tb')">前往淘宝支付</button>
                </div>
              </el-col>
            </el-row>
          </div>
        </div>
      </template>
    </el-dialog>

    <!-- 提示card模态框 -->
    <div v-if="showCardModal" class="modal">
      <div class="modal-content">
        <p style="font-size: 20px">未识别到合法卡密</p>
      </div>
    </div>

    <!-- 提示copy模态框 -->
    <div v-if="showCopyModal" class="modal">
      <div class="modal-content">
        <p style="font-size: 20px">复制成功</p>
      </div>
    </div>

    <!-- 提示submit card模态框 -->
    <div v-if="showSubmitCardModal" class="modal">
      <div class="modal-content">
        <p style="font-size: 20px">您已提交成功，等待客服核实</p>
      </div>
    </div>
    <!-- 提示submit card err模态框 -->
    <div v-if="showSubmitErrModal" class="modal">
      <div class="modal-content">
        <p style="font-size: 20px">{{ showSubmitErrInfo }}</p>
      </div>
    </div>

    <div v-show="notFoundVisible">
      <!-- 显示新的 div 的代码... -->
      <h1>订单不存在</h1>
    </div>
    <div v-show="finishedVisible">
      <!-- 显示新的 div 的代码... -->
      <h1>已付款成功</h1>
    </div>
    <div v-show="timeoutVisible">
      <!-- 显示新的 div 的代码... -->
      <h1>订单已超时</h1>
    </div>
  </div>
  <div v-show="exVisible">
    <!-- 显示新的 div 的代码... -->
    <h1>订单异常，请重新下单</h1>
  </div>
</template>
<script>
export default {
  name: 'Pay',
}
</script>
<script setup>
import {ElButton, ElMessage} from 'element-plus';
import {onMounted, ref, onUnmounted, onBeforeUnmount, watch, watchEffect} from 'vue';
import CountDown from 'vue-canvas-countdown';
import {cbExt, queryOrderSimple} from '@/api/payOrder';
import {useRoute} from 'vue-router';
import {WarningFilled} from '@element-plus/icons-vue';
import {formatTime} from "@/utils/format";
import QRCode from "qrcode";
import ClipboardJS from "clipboard";
import bgImage from '@/assets/od_info_bg.png'; // 背景图片

const backgroundImageStyle = `background-image: url(${bgImage});background-size: 100% 100%;`;
const backgroundYdImageStyle = `background-image: url(${bgImage});background-size: 100% 100%;border-radius: 10px;box-shadow: 0 4px 6px rgba(0, 0, 0, 0.5);top: -20px`;

// 弹窗控制标记
const dialogCountVisible = ref(true)
const payVisible = ref(false)
const payTypeVisible = ref(0)
const tmH5Visible = ref(false)
const copyInfoVisible = ref(false)
const readInfoVisible = ref(false)
const finishedVisible = ref(false)
const timeoutVisible = ref(false)
const exVisible = ref(false)
const notFoundVisible = ref(false)
const route = useRoute()

// ---------- 付款页 ----------------

const generateColor = (index) => {
  const hue = 220; // 蓝色的色调，范围为0-360
  const saturation = 90; // 蓝色的饱和度，范围为0-100
  const lightness = 100 - ((index - 0) * 8); // 蓝色的亮度，范围为0-100
  return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
};

// 直付扫码
const qrcodeUrl = ref('');

// ---------- 倒计时 ----------------
const fire = ref(0);
const tiping = {
  text: '倒计时进行中',
  color: '#fff'
};
const tipend = {
  text: '倒计时结束',
  color: '#fff'
};

const fireCD = async () => {
  // 配置参数（更多配置如下表）
  tiping.text = '匹配中';
  tiping.color = '#fff';
  tipend.text = '停止匹配';
  tipend.color = '#fff';

  // 启动倒计时(效果如上图所示)
  fire.value++;
};

const onStatusChange = async (payload) => {
  console.log('倒计时状态改变：', payload);
};

const onEnd = async () => {
  console.log('倒计时结束的回调函数');
};

// 识别card info
// --------------- card --------------------
// 输入的字符串
const inputString = ref('');
// 提取卡号和密码的正则表达式
const cardNumberRegex = /卡号[：:](\d{16})/;
const passwordRegex = /密码[：:](\d{16})/;

const card6000NumberRegex = /(\S{4}-\S{4}-\S{4}-\S{4})/;

// 使用 ref 来存储匹配到的卡号和密码
const card1100Number = ref(null);
const card6000Number = ref(null);
const password1100 = ref(null);

// 监听输入的字符串变化，进行匹配
watchEffect(() => {
  // 重置匹配结果
  card1100Number.value = null;
  card6000Number.value = null;
  password1100.value = null;

  // 匹配卡号
  const matchCardNumber = inputString.value.match(cardNumberRegex);
  if (matchCardNumber) {
    card1100Number.value = matchCardNumber[1];
  }

  // 匹配卡号
  const matchCard6000Number = inputString.value.match(card6000NumberRegex);
  if (matchCard6000Number) {
    card6000Number.value = matchCard6000Number[0];
  }

  // 匹配密码
  const matchPassword = inputString.value.match(passwordRegex);
  if (matchPassword) {
    password1100.value = matchPassword[1];
  }
});
// --------------- card --------------------

// 复制
const copyInfo = async () => {
  let cid = Number(payData.value.channel_code);
  let copyInfo = ''
  if (cid === 1006) {
    copyInfo = `【支付步骤】
-----
【如有历史账号记录，请先清除！】
【选择指定金额，核对粘贴指定账号再支付！】

1、指定金额：${payData.value.money}元
2、复制账号：${payData.value.account}
3、点击跳转下面地址：

 #小程序://腾讯充值/5I6C76SFkEcMeIj
 `
  } else {
    copyInfo = `${payData.value.account}`
  }
  console.log("copyInfo", copyInfo)
  const clipboard = new ClipboardJS('.btn-copy', {
    text: () => copyInfo
  });

  clipboard.on('success', () => {
    showCopyModal.value = true
    setTimeout(() => {
      showCopyModal.value = false
    }, 300)
    copyInfoVisible.value = true
    setTimeout(() => {
      copyInfoVisible.value = false
    }, 2000)
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });

  clipboard.on('error', () => {
    ElMessage({
      type: 'error',
      message: '复制异常'
    })
    clipboard.destroy(); // 销毁 ClipboardJS 实例
  });

};

const dialogYd1000Visible = ref(false)
const dialogYd1100Visible = ref(false)
const dialogYd1200Visible = ref(false)
const dialogYd2000Visible = ref(false)
const dialogYd3000Visible = ref(false)
const dialogYd4000Visible = ref(false)
const dialogYd5000Visible = ref(false)
const dialogYd6000Visible = ref(false)

const closeYdDialog = async () => {
  dialogYd1000Visible.value = false
  dialogYd1100Visible.value = false
  dialogYd1200Visible.value = false
  dialogYd2000Visible.value = false
  dialogYd3000Visible.value = false
  dialogYd4000Visible.value = false
  dialogYd5000Visible.value = false
  dialogYd6000Visible.value = false
}

const openYdVisible = async () => {
  let cid = payData.value.channel_code;
  if (cid >= 3000 && cid < 3099) {
    dialogYd3000Visible.value = true
  } else if (cid >= 2000 && cid < 2099) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd2000Visible.value = true
  } else if (cid >= 1100 && cid < 1199) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd1100Visible.value = true
  } else if (cid >= 1000 && cid < 1099) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd1000Visible.value = true
  }  else if (cid >= 1200 && cid < 1299) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd1200Visible.value = true
  } else if (cid >= 4000 && cid < 4099) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd4000Visible.value = true
  }else if (cid >= 6000 && cid < 6099) {
    startCountdown()
    readInfoVisible.value = true
    setTimeout(() => {
      readInfoVisible.value = false
    }, 3000)
    dialogYd6000Visible.value = true
  } else {

  }
}

const openPayHref = async (chan) => {
  let url = payData.value.resource_url;
  let money = Number(payData.value.money)
  if(chan === 'jd'){
    if(money === 10) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221446017%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 50){
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221107851%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 100) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221107845%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 200) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221107847%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 300) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221107846%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 500) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221107843%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 600) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221962859%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 800) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221107833%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 1000) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%221107842%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 2000) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%223348254%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 3000) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%223522645%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }else if (money === 5000) {
      url = 'openapp.jdmobile://virtual?params=%7B%22category%22%3A%22jump%22%2C%22des%22%3A%22productDetail%22%2C%22skuId%22%3A%223020581%22%2C%22sourceType%22%3A%22JSHOP_SOURCE_TYPE%22%2C%22sourceValue%22%3A%22JSHOP_SOURCE_VALUE%22%7D'
    }
  }
  if(chan === 'tb'){
    if(money === 10) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A110&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 50){
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A150&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 100) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A1100&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 200) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A1200&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 300) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A1300&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 500) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A1500&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 600) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A1600&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 800) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A1800&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 1000) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A11000&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 2000) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A12000&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 3000) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A13000&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }else if (money === 5000) {
      url = 'tbopen://m.taobao.com/tbopen/index.html?h5Url=https://main.m.taobao.com/search/index.html?q=%25E4%25BA%25AC%25E4%25B8%259Ce%25E5%258D%25A15000&action=ali.open.nav&module=h5&bootImage=0&slk_sid=&slk_t=&slk_gid=gid_er_er%7Cgid_er_af_pop&afcPromotionOpen=false&source=slk_dp'
    }
  }

  // 等待200毫秒后跳转
  setTimeout(() => {
    // window.location.href = payData.value.resource_url;
    window.open(url, '_blank')
  }, 100);
}
const openPay = async () => {
  let cid = Number(payData.value.channel_code);
  let copyInfo = ''
  if (cid === 1006) {
    copyInfo = `【支付步骤】
-----
【如有历史账号记录，请先清除！】
【选择指定金额，核对粘贴指定账号再支付！】

1、指定金额：${payData.value.money}元
2、复制账号：${payData.value.account}
3、点击跳转下面地址：

 #小程序://腾讯充值/5I6C76SFkEcMeIj
 `
    payData.value.resource_url = "weixin://"
  } else {
    copyInfo = `${payData.value.account}`
  }
  console.log("copyInfo", copyInfo)

  const clipboardX = new ClipboardJS('.btn-copy', {
    text: () => copyInfo
  });

  clipboardX.on('success', () => {
    console.log('复制成功', copyInfo)

    // 等待200毫秒后跳转
    setTimeout(() => {
      // window.location.href = payData.value.resource_url;
      window.open(payData.value.resource_url, '_blank')
    }, 100);

    clipboardX.destroy(); // 销毁 ClipboardJS 实例
  });

  clipboardX.on('error', () => {
    ElMessage({
      type: 'error',
      message: '复制异常'
    })
    console.log('复制异常', copyInfo)

    clipboardX.destroy(); // 销毁 ClipboardJS 实例
  });

};

// 添加一个空变量作为定时器的 ID
let timerId = null;
let timerExp = null;
let timerYD = null;

onMounted(() => {
  // 启动倒计时
  fireCD();
  // 启动定时器，每秒钟请求一次 HTTP 接口
  // timerId = setInterval(queryOrder, 1000);
  startCountdownQryOrder();

  // setCacheControl(180) // 设置缓存时间为180秒
});

onUnmounted(() => {
  // 组件销毁时清除定时器
  clearInterval(timerId);
});

const payData = ref({
  money: 0,
  exp_time: 0,
  channel_code: 0,
  account: '',
  order_id: '',
  resource_url: '',
  status: 0,
  ext: '',
})

const reqCnt = ref(0)

const queryOrder = async (timerId) => {
  try {
    const orderId = route.query.orderId;
    console.log(orderId)
    const result = await queryOrderSimple({order_id: orderId}); // 发送 HTTP 请求
    let nowTime = new Date().getTime();
    let resExp = new Date(result.data?.exp_time).getTime();
    console.log(nowTime)
    console.log(resExp)
    payData.value = result.data
    console.log(payData.value)
    const content = result.data.resource_url;
    const account = result.data.account;
    if (content && account) {
      console.log("qry time id ", timerId)
      clearInterval(timerId);
      startCountdownExp();
      // 如果状态发生变化，则停止定时器
      dialogCountVisible.value = false;

      if (result.code === 7) {
        dialogCountVisible.value = false;
        notFoundVisible.value = true;
      } else if (result.code === 0) {

        payTypeVisible.value = Number(result.data.channel_code);
        //如果 payData.value.resource_url 中包含 main.m.taobao.com
        if (content.includes('main.m.taobao.com')) {
          tmH5Visible.value = true;
        }
        if (content) {
          QRCode.toDataURL(content)
              .then((dataUrl) => {
                qrcodeUrl.value = dataUrl
              })
              .catch((error) => {
                console.error('Failed to generate QR code:', error);
              });
        } else {
          // 付款码异常
        }

        if (result.data.status === 1) {
          finishedVisible.value = true;
        }
        if (result.data.status === 2) {
          payVisible.value = true;
        }
        if (result.data.status === 3) {
          timeoutVisible.value = true;
        }
      }
    } else if (result.data?.status === 0) {
      dialogCountVisible.value = false;
      clearInterval(timerId);
      exVisible.value = true;
    } else if (result.data?.status === 2 && resExp < nowTime) {
      dialogCountVisible.value = false;
      clearInterval(timerId);
      exVisible.value = true;
      console.log('超时')
    } else if (reqCnt.value < 10) {
      reqCnt.value++
    } else {
      dialogCountVisible.value = false;
      clearInterval(timerId);
      exVisible.value = true;
    }
    // if (result) {
    //   clearInterval(timerId); // 如果状态发生变化，则停止定时器
    // }

  } catch (error) {
    console.log(error);
  }
}

// 倒计时数组
const countdowns = ref([]);

// 计算倒计时
// const calculateCountdown = () => {
//   setInterval(() => {
//     const currentTime = new Date();
//     const timeLimit = new Date(payData.value.exp_time);
//     const timeDiffInSeconds = (timeLimit - currentTime) / 1000;
//     countdowns.value[0] = timeDiffInSeconds > 0 ? Math.floor(timeDiffInSeconds) : -1;
//   }, 1000);
// };

// 引导倒计时
const countdownTime = ref(3);

const countdown = (countdownTimeRef, timerRef) => {
  countdownTimeRef.value--;

  // 监听倒计时时间变化，可以在这里执行倒计时结束后的操作
  watch(countdownTimeRef, (newVal) => {
    if (newVal === 0) {
      // 在此处执行倒计时结束后的操作
      console.log(`倒计时结束 - ${countdownTimeRef.value}`);
      countdownTime.value = 3;
      // 清除定时器
      clearInterval(timerRef);
    }
  });
};

const startCountdown = () => {
  timerYD = setInterval(() => countdown(countdownTime, timerYD), 1000);
};

const countdownTimeQryOrder = ref(20);
const countdownQryOrder = (countdownTimeRef, timerRef) => {
  countdownTimeRef.value--;

  // 监听倒计时时间变化，可以在这里执行倒计时结束后的操作
  watch(countdownTimeRef, (newVal) => {
    console.log(newVal)
    queryOrder(timerRef)
    if (newVal === 0) {
      // 在此处执行倒计时结束后的操作
      console.log(`倒计时结束 - ${countdownTimeRef.value}`);
      // 清除定时器
      clearInterval(timerRef);
    }
  });
};

const startCountdownQryOrder = () => {
  timerId = setInterval(() => countdownQryOrder(countdownTimeQryOrder, timerId), 1000);
};

const startCountdownExp = () => {
  timerExp = setInterval(() => {
    const currentTime = new Date();
    const timeLimit = new Date(payData.value.exp_time);
    const timeDiffInSeconds = (timeLimit - currentTime) / 1000;
    countdowns.value[0] = timeDiffInSeconds > 0 ? Math.floor(timeDiffInSeconds) : -1;
    // console.log('timeLimit', timeLimit);
    // console.log('timeDiffInSeconds', timeDiffInSeconds);
    // console.log(timerExp)
    if (timeDiffInSeconds < 0) {
      clearInterval(timerExp);
    }
  }, 1000);
};

// 在组件销毁前清除定时器
onBeforeUnmount(() => {
  clearInterval(timerYD);
});

// card info
const dialog1100CardVisible = ref(false)

const open1100CardVisible = async () => {
  dialog1100CardVisible.value = true
}
const close1100CardVisible = () => {
  dialog1100CardVisible.value = false
}

// card info
const dialog6000CardVisible = ref(false)

const open6000CardVisible = async () => {
  dialog6000CardVisible.value = true
}
const close6000CardVisible = () => {
  dialog6000CardVisible.value = false
}

//modal提示
// 控制是否显示模态框
const showCopyModal = ref(false);
const showSubmitCardModal = ref(false);
const showSubmitErrModal = ref(false);
const showSubmitErrInfo = ref();
const showCardModal = ref(false);
const submit1100CardInfo = async () => {
  let c = String(card1100Number.value)
  let p = String(password1100.value)
  if (c && p) {
    payData.value.ext = c + "_" + p
    const cbRes = await cbExt({...payData.value})
    await close1100CardVisible()
    if (cbRes.code === 0) {
      showSubmitCardModal.value = true;
      // 设置一段时间后隐藏模态框（例如，3秒后隐藏）
      setTimeout(() => {
        showSubmitCardModal.value = false;
      }, 2000);
    } else if (cbRes.code === 7) {
      showSubmitErrInfo.value = cbRes.msg;
      showSubmitErrModal.value = true;
      // 设置一段时间后隐藏模态框（例如，3秒后隐藏）
      setTimeout(() => {
        showSubmitErrModal.value = false;
      }, 2000);
    }
  } else {
    await warnCardInfo();
  }
}

const submit6000CardInfo = async () => {
  let c = String(card6000Number.value)
  if (c) {
    payData.value.ext = c
    const cbRes = await cbExt({...payData.value})
    await close6000CardVisible()
    if (cbRes.code === 0) {
      showSubmitCardModal.value = true;
      // 设置一段时间后隐藏模态框（例如，3秒后隐藏）
      setTimeout(() => {
        showSubmitCardModal.value = false;
      }, 2000);
    } else if (cbRes.code === 7) {
      showSubmitErrInfo.value = cbRes.msg;
      showSubmitErrModal.value = true;
      // 设置一段时间后隐藏模态框（例如，3秒后隐藏）
      setTimeout(() => {
        showSubmitErrModal.value = false;
      }, 2000);
    }
  } else {
    await warnCardInfo();
  }
}

const warnCardInfo = () => {
  // 显示模态框
  showCardModal.value = true;

  // 设置一段时间后隐藏模态框（例如，3秒后隐藏）
  setTimeout(() => {
    showCardModal.value = false;
  }, 1000);
}


</script>

<style scoped>
.c_container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.c_content {
  text-align: center;
  color: #333;
}

h1 {
  font-size: 36px;
  margin-bottom: 24px;
}


.c_button {
  padding: 12px 24px;
  font-size: 18px;
  margin-top: 6px;
  width: 80%;
  height: 42px;
}

.p_container {
  display: flex;
  justify-content: space-between;
  height: 45vh;
}

.p_blue-section {
  flex-grow: 1;
}

.p_content {
  text-align: center;
  color: #333;
  position: absolute;
  top: 50px;
  left: 5%;
  right: 5%;
  height: 280px;
  background-color: #f2f2f2;
  border-radius: 20px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_inner {
  text-align: center;
  color: #333;
  position: absolute;
  top: 320px;
  left: 5%;
  right: 5%;
  height: 200px;
  background-color: #f2f2f2;
  border-radius: 10px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_yd_inner {
  text-align: center;
  color: #333;
  top: 320px;
  left: 5%;
  right: 5%;
  height: 200px;
  background-color: #f2f2f2;
  border-radius: 10px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_button {
  text-align: center;
  color: #333;
  position: absolute;
  top: 550px;
  left: 5%;
  right: 5%;
}

.p_content_button_qr {
  text-align: center;
  color: #333;
  position: absolute;
  top: 580px;
  left: 5%;
  right: 5%;
}

.p_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(to right, #064954, #125280, #1247c9);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.yd_p_content_button_qr {
  text-align: center;
  color: #333;
  position: absolute;
  left: 5%;
  right: 5%;
}

.yd_read_p_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(to right, #a5abb4, rgba(122, 129, 140, 0.99), #a5abb4);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.yd_p_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(90deg, #5498ff 1%, #00d9d0 100%);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.yd_p_jd_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(90deg, #d70b0b 1%, #d9003a 100%);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.yd_p_tb_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(90deg, #d77639 1%, rgba(210, 80, 41, 0.99) 100%);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.copy_button {
  border: none;
  font-size: 16px;
  color: #e7dfdf;
  background: linear-gradient(to right, #d71010, #ab4f34);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 30px;
  margin-top: 6px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
}

.copy_success_button {
  border: none;
  font-size: 16px;
  color: #e7dfdf;
  background: linear-gradient(to right, #0d650a, #71ab34);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 30px;
  margin-top: 6px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
}

.copy-tip {
  display: flex;
  justify-content: space-around;
  align-items: center;
  background: linear-gradient(to right, #d71010, #064954);
  color: white;
  font-size: 14px;
  height: 24px;
  margin-top: 6px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
}

.jtone {
  position: relative;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid white;
  border-bottom: 12px solid transparent;
}

.jtone::before {
  content: "";
  position: absolute;
  top: -12px;
  left: -13px;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid #9f330c;
  border-bottom: 12px solid transparent;
}

.jttwo {
  position: relative;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid white;
  border-bottom: 12px solid transparent;
}

.jttwo::before {
  content: "";
  position: absolute;
  top: -12px;
  left: -13px;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid rgba(126, 60, 46, 0.99);
  border-bottom: 12px solid transparent;
}

.jtthree {
  position: relative;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid white;
  border-bottom: 12px solid transparent;
}

.jtthree::before {
  content: "";
  position: absolute;
  top: -12px;
  left: -13px;
  width: 0;
  height: 0;
  border-top: 12px solid transparent;
  border-left: 10px solid rgb(77, 67, 61);
  border-bottom: 12px solid transparent;
}

.medicine-money-bag {
  background: rgba(215, 197, 197, 0.1);
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  height: 30px;
}

.medicine-bag {
  background: rgba(215, 197, 197, 0.1);
  border: 1px dashed rgba(59, 28, 23, 0.99);
  margin-top: 6px;
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  font-size: 18px;
  height: 30px;
}

.medicine-bag-qr {
  background: rgba(215, 197, 197, 0.1);
  border: 1px dashed rgba(59, 28, 23, 0.99);
  border-radius: 5px;
  padding-top: 5px;
  margin-left: 10%;
  margin-right: 10%;
  width: 80%;
  font-size: 18px;
  height: 180px;
}

.p_content_card_info_inner {
  text-align: center;
  color: #333;
  position: absolute;
  top: 380px;
  left: 5%;
  right: 5%;
  height: 270px;
  background-color: #f2f2f2;
  border-radius: 10px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.p_content_card_info_button {
  text-align: center;
  color: #333;
  position: absolute;
  top: 340px;
  left: 5%;
  right: 5%;
}

.p_content_card_submit_button {
  text-align: center;
  color: #333;
  position: absolute;
  top: 240px;
  left: 5%;
  right: 5%;
}

.medicine-jw-bag {
  background: rgba(215, 197, 197, 0.1);
  border: 1px solid rgba(59, 28, 23, 0.99);
  margin-top: 6px;
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  font-size: 18px;
  height: 100px;
}

.medicine-ec-bag {
  background: rgba(215, 197, 197, 0.1);
  border: 1px solid rgba(59, 28, 23, 0.99);
  margin-top: 6px;
  text-align: center;
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  font-size: 18px;
  height: 60px;
}

.medicine-jw-card-info-bag {
  background: rgba(220, 200, 200, 0.1);
  border: 2px dashed rgba(59, 28, 23, 0.5);
  margin-top: 8px;
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  font-size: 18px;
  height: 50px;
}

.medicine-jw-card-info-submit-bag {
  background: rgba(220, 200, 200, 0.1);
  border: 1px dashed rgba(59, 28, 23, 0.5);
  margin-top: 8px;
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  font-size: 18px;
  height: 120px;
}
.medicine-ec-card-info-submit-bag {
  background: rgba(220, 200, 200, 0.1);
  border: 1px dashed rgba(59, 28, 23, 0.5);
  margin-top: 8px;
  border-radius: 5px;
  padding-top: 10px;
  margin-left: 5%;
  margin-right: 5%;
  width: 90%;
  font-size: 18px;
  height: 80px;
}

textarea {
  resize: none;
  border: none;
  background-color: transparent;
  width: 100%;
  height: 100%;
  font-size: 16px;
  line-height: 1.5;
  outline: none;
  overflow: auto;
}

/* 设置 placeholder 文字样式 */
textarea::placeholder {
  color: #c4bdbd; /* 设置灰色提示文字颜色 */
}

.result-container {
  border-radius: 5px;
}

.result-container p {
  margin: 5px 0;
}

.p_submit_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(to right, #e8731f, #d54d11, #c91258);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.p_submit_success_button {
  border: none;
  padding: 12px 24px;
  font-size: 22px;
  color: #e7dfdf;
  background: linear-gradient(to right, #064954, #2c9a12, #075cbd);
  margin-top: 6px;
  width: 80%;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  height: 50px;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* 透明黑色背景 */
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-content {
  background-color: #111111;
  color: #FFFFFF;
  padding: 20px;
  border-radius: 5px;
  text-align: center;
}

.modal-content p {
  margin: 0;
}
</style>
